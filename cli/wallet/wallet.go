package wallet

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/nspcc-dev/neo-go/cli/cmdargs"
	"github.com/nspcc-dev/neo-go/cli/flags"
	"github.com/nspcc-dev/neo-go/cli/input"
	"github.com/nspcc-dev/neo-go/cli/options"
	"github.com/nspcc-dev/neo-go/cli/txctx"
	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/encoding/address"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/neo"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/manifest"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/vm"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"
)

const (
	// EnterPasswordPrompt is a prompt used to ask the user for a password.
	EnterPasswordPrompt = "Enter password > "
	// EnterNewPasswordPrompt is a prompt used to ask the user for a password on
	// account creation.
	EnterNewPasswordPrompt = "Enter new password > "
	// EnterOldPasswordPrompt is a prompt used to ask the user for an old password.
	EnterOldPasswordPrompt = "Enter old password > "
	// ConfirmPasswordPrompt is a prompt used to confirm the password.
	ConfirmPasswordPrompt = "Confirm password > "
)

var (
	errNoPath                 = errors.New("wallet path is mandatory and should be passed using (--wallet, -w) flags or via wallet config using --wallet-config flag")
	errConflictingWalletFlags = errors.New("--wallet flag conflicts with --wallet-config flag, please, provide one of them to specify wallet location")
	errPhraseMismatch         = errors.New("the entered pass-phrases do not match. Maybe you have misspelled them")
	errNoStdin                = errors.New("can't read wallet from stdin for this command")
)

var (
	walletPathFlag = cli.StringFlag{
		Name:  "wallet, w",
		Usage: "Target location of the wallet file ('-' to read from stdin); conflicts with --wallet-config flag.",
	}
	walletConfigFlag = cli.StringFlag{
		Name:  "wallet-config",
		Usage: "Target location of the wallet config file; conflicts with --wallet flag.",
	}
	wifFlag = cli.StringFlag{
		Name:  "wif",
		Usage: "WIF to import",
	}
	decryptFlag = cli.BoolFlag{
		Name:  "decrypt, d",
		Usage: "Decrypt encrypted keys.",
	}
	inFlag = cli.StringFlag{
		Name:  "in",
		Usage: "file with JSON transaction",
	}
	fromAddrFlag = flags.AddressFlag{
		Name:  "from",
		Usage: "Address to send an asset from",
	}
	toAddrFlag = flags.AddressFlag{
		Name:  "to",
		Usage: "Address to send an asset to",
	}
)

// NewCommands returns 'wallet' command.
func NewCommands() []cli.Command {
	claimFlags := []cli.Flag{
		walletPathFlag,
		walletConfigFlag,
		txctx.GasFlag,
		txctx.SysGasFlag,
		txctx.OutFlag,
		txctx.ForceFlag,
		flags.AddressFlag{
			Name:  "address, a",
			Usage: "Address to claim GAS for",
		},
	}
	claimFlags = append(claimFlags, options.RPC...)
	signFlags := []cli.Flag{
		walletPathFlag,
		walletConfigFlag,
		txctx.OutFlag,
		inFlag,
		flags.AddressFlag{
			Name:  "address, a",
			Usage: "Address to use",
		},
	}
	signFlags = append(signFlags, options.RPC...)
	return []cli.Command{{
		Name:  "wallet",
		Usage: "create, open and manage a NEO wallet",
		Subcommands: []cli.Command{
			{
				Name:      "claim",
				Usage:     "claim GAS",
				UsageText: "neo-go wallet claim -w wallet [--wallet-config path] [-g gas] [-e sysgas] -a address -r endpoint [-s timeout] [--out file] [--force]",
				Action:    claimGas,
				Flags:     claimFlags,
			},
			{
				Name:      "init",
				Usage:     "create a new wallet",
				UsageText: "neo-go wallet init -w wallet [--wallet-config path] [-a]",
				Action:    createWallet,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					cli.BoolFlag{
						Name:  "account, a",
						Usage: "Create a new account",
					},
				},
			},
			{
				Name:      "change-password",
				Usage:     "change password for accounts",
				UsageText: "neo-go wallet change-password -w wallet -a address",
				Action:    changePassword,
				Flags: []cli.Flag{
					walletPathFlag,
					flags.AddressFlag{
						Name:  "address, a",
						Usage: "address to change password for",
					},
				},
			},
			{
				Name:      "convert",
				Usage:     "convert addresses from existing NEO2 NEP6-wallet to NEO3 format",
				UsageText: "neo-go wallet convert -w legacywallet [--wallet-config path] -o n3wallet",
				Action:    convertWallet,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					cli.StringFlag{
						Name:  "out, o",
						Usage: "where to write converted wallet",
					},
				},
			},
			{
				Name:      "create",
				Usage:     "add an account to the existing wallet",
				UsageText: "neo-go wallet create -w wallet [--wallet-config path]",
				Action:    addAccount,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
				},
			},
			{
				Name:      "dump",
				Usage:     "check and dump an existing NEO wallet",
				UsageText: "neo-go wallet dump -w wallet [--wallet-config path] [-d]",
				Description: `Prints the given wallet (via -w option or via wallet configuration file) in JSON
   format to the standard output. If -d is given, private keys are unencrypted and
   displayed in clear text on the console! Be very careful with this option and
   don't use it unless you know what you're doing.
`,
				Action: dumpWallet,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					decryptFlag,
				},
			},
			{
				Name:      "dump-keys",
				Usage:     "dump public keys for account",
				UsageText: "neo-go wallet dump-keys -w wallet [--wallet-config path] [-a address]",
				Action:    dumpKeys,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					flags.AddressFlag{
						Name:  "address, a",
						Usage: "address to print public keys for",
					},
				},
			},
			{
				Name:      "export",
				Usage:     "export keys for address",
				UsageText: "export -w wallet [--wallet-config path] [--decrypt] [<address>]",
				Description: `Prints the key for the given account to the standard output. It uses NEP-2
   encrypted format by default (the way NEP-6 wallets store it) or WIF format if
   -d option is given. In the latter case the key can be displayed in clear text
   on the console, so be extremely careful with this option and don't use unless
   you really need it and know what you're doing.
`,
				Action: exportKeys,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					decryptFlag,
				},
			},
			{
				Name:      "import",
				Usage:     "import WIF of a standard signature contract",
				UsageText: "import -w wallet [--wallet-config path] --wif <wif> [--name <account_name>]",
				Action:    importWallet,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					wifFlag,
					cli.StringFlag{
						Name:  "name, n",
						Usage: "Optional account name",
					},
					cli.StringFlag{
						Name:  "contract",
						Usage: "Verification script for custom contracts",
					},
				},
			},
			{
				Name:  "import-multisig",
				Usage: "import multisig contract",
				UsageText: "import-multisig -w wallet [--wallet-config path] --wif <wif> [--name <account_name>] --min <n>" +
					" [<pubkey1> [<pubkey2> [...]]]",
				Action: importMultisig,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					wifFlag,
					cli.StringFlag{
						Name:  "name, n",
						Usage: "Optional account name",
					},
					cli.IntFlag{
						Name:  "min, m",
						Usage: "Minimal number of signatures",
					},
				},
			},
			{
				Name:      "import-deployed",
				Usage:     "import deployed contract",
				UsageText: "import-deployed -w wallet [--wallet-config path] --wif <wif> --contract <hash> [--name <account_name>]",
				Action:    importDeployed,
				Flags: append([]cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					wifFlag,
					cli.StringFlag{
						Name:  "name, n",
						Usage: "Optional account name",
					},
					flags.AddressFlag{
						Name:  "contract, c",
						Usage: "Contract hash or address",
					},
				}, options.RPC...),
			},
			{
				Name:      "remove",
				Usage:     "remove an account from the wallet",
				UsageText: "remove -w wallet [--wallet-config path] [--force] --address <addr>",
				Action:    removeAccount,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					txctx.ForceFlag,
					flags.AddressFlag{
						Name:  "address, a",
						Usage: "Account address or hash in LE form to be removed",
					},
				},
			},
			{
				Name:      "sign",
				Usage:     "cosign transaction with multisig/contract/additional account",
				UsageText: "sign -w wallet [--wallet-config path] --address <address> --in <file.in> [--out <file.out>] [-r <endpoint>]",
				Description: `Signs the given (in file.in) context (which must be a transaction
   signing context) for the given address using the given wallet. This command can
   output the resulting JSON (with additional signature added) right to the console
   (if no file.out and no RPC endpoint specified) or into a file (which can be the
   same as input one). If an RPC endpoint is given it'll also try to construct a
   complete transaction and send it via RPC (printing its hash if everything is OK).
`,
				Action: signStoredTransaction,
				Flags:  signFlags,
			},
			{
				Name:      "strip-keys",
				Usage:     "remove private keys for all accounts",
				UsageText: "neo-go wallet strip-keys -w wallet [--wallet-config path] [--force]",
				Description: `Removes private keys for all accounts from the given wallet. Notice,
   this is a very dangerous action (you can lose keys if you don't have a wallet
   backup) that should not be performed unless you know what you're doing. It's
   mostly useful for creation of special wallets that can be used to create
   transactions, but can't be used to sign them (offline signing).
`,
				Action: stripKeys,
				Flags: []cli.Flag{
					walletPathFlag,
					walletConfigFlag,
					txctx.ForceFlag,
				},
			},
			{
				Name:        "nep17",
				Usage:       "work with NEP-17 contracts",
				Subcommands: newNEP17Commands(),
			},
			{
				Name:        "nep11",
				Usage:       "work with NEP-11 contracts",
				Subcommands: newNEP11Commands(),
			},
			{
				Name:        "candidate",
				Usage:       "work with candidates",
				Subcommands: newValidatorCommands(),
			},
		},
	}}
}

func claimGas(ctx *cli.Context) error {
	return handleNeoAction(ctx, func(contract *neo.Contract, shash util.Uint160, _ *wallet.Account) (*transaction.Transaction, error) {
		return contract.TransferUnsigned(shash, shash, big.NewInt(0), nil)
	})
}

func changePassword(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, _, err := openWallet(ctx, false)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()
	if len(wall.Accounts) == 0 {
		return cli.NewExitError("wallet has no accounts", 1)
	}
	addrFlag := ctx.Generic("address").(*flags.Address)
	if addrFlag.IsSet {
		// Check for account presence first before asking for password.
		acc := wall.GetAccount(addrFlag.Uint160())
		if acc == nil {
			return cli.NewExitError("account is missing", 1)
		}
	}

	oldPass, err := input.ReadPassword(EnterOldPasswordPrompt)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("Error reading old password: %w", err), 1)
	}

	for i := range wall.Accounts {
		if addrFlag.IsSet && wall.Accounts[i].Address != addrFlag.String() {
			continue
		}
		err := wall.Accounts[i].Decrypt(oldPass, wall.Scrypt)
		if err != nil {
			return cli.NewExitError(fmt.Errorf("unable to decrypt account %s: %w", wall.Accounts[i].Address, err), 1)
		}
	}

	pass, err := readNewPassword()
	if err != nil {
		return cli.NewExitError(fmt.Errorf("Error reading new password: %w", err), 1)
	}
	for i := range wall.Accounts {
		if addrFlag.IsSet && wall.Accounts[i].Address != addrFlag.String() {
			continue
		}
		err := wall.Accounts[i].Encrypt(pass, wall.Scrypt)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
	}
	err = wall.Save()
	if err != nil {
		return cli.NewExitError(fmt.Errorf("Error saving the wallet: %w", err), 1)
	}
	return nil
}

func convertWallet(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, pass, err := newWalletV2FromFile(ctx.String("wallet"), ctx.String("wallet-config"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	out := ctx.String("out")
	if len(out) == 0 {
		return cli.NewExitError("missing out path", 1)
	}
	newWallet, err := wallet.NewWallet(out)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	newWallet.Scrypt = wall.Scrypt

	for _, acc := range wall.Accounts {
		if len(wall.Accounts) != 1 || pass == nil {
			password, err := input.ReadPassword(fmt.Sprintf("Enter password for account %s (label '%s') > ", acc.Address, acc.Label))
			if err != nil {
				return cli.NewExitError(fmt.Errorf("Error reading password: %w", err), 1)
			}
			pass = &password
		}

		newAcc, err := acc.convert(*pass, wall.Scrypt)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		newWallet.AddAccount(newAcc)
	}
	if err := newWallet.Save(); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func addAccount(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, pass, err := openWallet(ctx, true)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()

	if err := createAccount(wall, pass); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func exportKeys(ctx *cli.Context) error {
	wall, pass, err := readWallet(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()

	var addr string

	decrypt := ctx.Bool("decrypt")
	if ctx.NArg() == 0 && decrypt {
		return cli.NewExitError(errors.New("address must be provided if '--decrypt' flag is used"), 1)
	} else if ctx.NArg() > 0 {
		// check address format just to catch possible typos
		addr = ctx.Args().First()
		_, err := address.StringToUint160(addr)
		if err != nil {
			return cli.NewExitError(fmt.Errorf("can't parse address: %w", err), 1)
		}
	}

	var wifs []string

loop:
	for _, a := range wall.Accounts {
		if addr != "" && a.Address != addr {
			continue
		}

		for i := range wifs {
			if a.EncryptedWIF == wifs[i] {
				continue loop
			}
		}

		wifs = append(wifs, a.EncryptedWIF)
	}

	for _, wif := range wifs {
		if decrypt {
			if pass == nil {
				password, err := input.ReadPassword(EnterPasswordPrompt)
				if err != nil {
					return cli.NewExitError(fmt.Errorf("Error reading password: %w", err), 1)
				}
				pass = &password
			}

			pk, err := keys.NEP2Decrypt(wif, *pass, wall.Scrypt)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			wif = pk.WIF()
		}

		fmt.Fprintln(ctx.App.Writer, wif)
	}

	return nil
}

func importMultisig(ctx *cli.Context) error {
	wall, _, err := openWallet(ctx, true)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()

	m := ctx.Int("min")
	if ctx.NArg() < m {
		return cli.NewExitError(errors.New("insufficient number of public keys"), 1)
	}

	args := []string(ctx.Args())
	pubs := make([]*keys.PublicKey, len(args))

	for i := range args {
		pubs[i], err = keys.NewPublicKeyFromString(args[i])
		if err != nil {
			return cli.NewExitError(fmt.Errorf("can't decode public key %d: %w", i, err), 1)
		}
	}

	acc, err := newAccountFromWIF(ctx.App.Writer, ctx.String("wif"), wall.Scrypt)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err := acc.ConvertMultisig(m, pubs); err != nil {
		return cli.NewExitError(err, 1)
	}

	if acc.Label == "" {
		acc.Label = ctx.String("name")
	}
	if err := addAccountAndSave(wall, acc); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func importDeployed(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, _, err := openWallet(ctx, true)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()

	rawHash := ctx.Generic("contract").(*flags.Address)
	if !rawHash.IsSet {
		return cli.NewExitError("contract hash was not provided", 1)
	}

	acc, err := newAccountFromWIF(ctx.App.Writer, ctx.String("wif"), wall.Scrypt)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	gctx, cancel := options.GetTimeoutContext(ctx)
	defer cancel()

	c, err := options.GetRPCClient(gctx, ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	cs, err := c.GetContractStateByHash(rawHash.Uint160())
	if err != nil {
		return cli.NewExitError(fmt.Errorf("can't fetch contract info: %w", err), 1)
	}
	md := cs.Manifest.ABI.GetMethod(manifest.MethodVerify, -1)
	if md == nil || md.ReturnType != smartcontract.BoolType {
		return cli.NewExitError("contract has no `verify` method with boolean return", 1)
	}
	acc.Address = address.Uint160ToString(cs.Hash)
	acc.Contract.Script = cs.NEF.Script
	acc.Contract.Parameters = acc.Contract.Parameters[:0]
	for _, p := range md.Parameters {
		acc.Contract.Parameters = append(acc.Contract.Parameters, wallet.ContractParam{
			Name: p.Name,
			Type: p.Type,
		})
	}
	acc.Contract.Deployed = true

	if acc.Label == "" {
		acc.Label = ctx.String("name")
	}
	if err := addAccountAndSave(wall, acc); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func importWallet(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, _, err := openWallet(ctx, true)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()

	acc, err := newAccountFromWIF(ctx.App.Writer, ctx.String("wif"), wall.Scrypt)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if ctrFlag := ctx.String("contract"); ctrFlag != "" {
		ctr, err := hex.DecodeString(ctrFlag)
		if err != nil {
			return cli.NewExitError("invalid contract", 1)
		}
		acc.Contract.Script = ctr
	}

	if acc.Label == "" {
		acc.Label = ctx.String("name")
	}
	if err := addAccountAndSave(wall, acc); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func removeAccount(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, _, err := openWallet(ctx, true)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()

	addr := ctx.Generic("address").(*flags.Address)
	if !addr.IsSet {
		return cli.NewExitError("valid account address must be provided", 1)
	}
	acc := wall.GetAccount(addr.Uint160())
	if acc == nil {
		return cli.NewExitError("account wasn't found", 1)
	}

	if !ctx.Bool("force") {
		fmt.Fprintf(ctx.App.Writer, "Account %s will be removed. This action is irreversible.\n", addr.Uint160())
		if ok := askForConsent(ctx.App.Writer); !ok {
			return nil
		}
	}

	if err := wall.RemoveAccount(acc.Address); err != nil {
		return cli.NewExitError(fmt.Errorf("error on remove: %w", err), 1)
	}
	if err := wall.Save(); err != nil {
		return cli.NewExitError(fmt.Errorf("error while saving wallet: %w", err), 1)
	}
	return nil
}

func askForConsent(w io.Writer) bool {
	response, err := input.ReadLine("Are you sure? [y/N]: ")
	if err == nil {
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		}
	}
	fmt.Fprintln(w, "Cancelled.")
	return false
}

func dumpWallet(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, pass, err := readWallet(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()
	if ctx.Bool("decrypt") {
		if pass == nil {
			password, err := input.ReadPassword(EnterPasswordPrompt)
			if err != nil {
				return cli.NewExitError(fmt.Errorf("Error reading password: %w", err), 1)
			}
			pass = &password
		}
		for i := range wall.Accounts {
			// Just testing the decryption here.
			err := wall.Accounts[i].Decrypt(*pass, wall.Scrypt)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
		}
	}
	fmtPrintWallet(ctx.App.Writer, wall)
	return nil
}

func dumpKeys(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, _, err := readWallet(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()
	accounts := wall.Accounts

	addrFlag := ctx.Generic("address").(*flags.Address)
	if addrFlag.IsSet {
		acc := wall.GetAccount(addrFlag.Uint160())
		if acc == nil {
			return cli.NewExitError("account is missing", 1)
		}
		accounts = []*wallet.Account{acc}
	}

	hasPrinted := false
	for _, acc := range accounts {
		pub, ok := vm.ParseSignatureContract(acc.Contract.Script)
		if ok {
			if hasPrinted {
				fmt.Fprintln(ctx.App.Writer)
			}
			fmt.Fprintf(ctx.App.Writer, "%s (simple signature contract):\n", acc.Address)
			fmt.Fprintln(ctx.App.Writer, hex.EncodeToString(pub))
			hasPrinted = true
			continue
		}
		n, bs, ok := vm.ParseMultiSigContract(acc.Contract.Script)
		if ok {
			if hasPrinted {
				fmt.Fprintln(ctx.App.Writer)
			}
			fmt.Fprintf(ctx.App.Writer, "%s (%d out of %d multisig contract):\n", acc.Address, n, len(bs))
			for i := range bs {
				fmt.Fprintln(ctx.App.Writer, hex.EncodeToString(bs[i]))
			}
			hasPrinted = true
			continue
		}
		if addrFlag.IsSet {
			return cli.NewExitError(fmt.Errorf("unknown script type for address %s", address.Uint160ToString(addrFlag.Uint160())), 1)
		}
	}
	return nil
}

func stripKeys(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	wall, _, err := readWallet(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	defer wall.Close()
	if !ctx.Bool("force") {
		fmt.Fprintln(ctx.App.Writer, "All private keys for all accounts will be removed from the wallet. This action is irreversible.")
		if ok := askForConsent(ctx.App.Writer); !ok {
			return nil
		}
	}
	for _, a := range wall.Accounts {
		a.EncryptedWIF = ""
	}
	if err := wall.Save(); err != nil {
		return cli.NewExitError(fmt.Errorf("error while saving wallet: %w", err), 1)
	}
	return nil
}

func createWallet(ctx *cli.Context) error {
	if err := cmdargs.EnsureNone(ctx); err != nil {
		return err
	}
	path := ctx.String("wallet")
	configPath := ctx.String("wallet-config")

	if len(path) != 0 && len(configPath) != 0 {
		return errConflictingWalletFlags
	}
	if len(path) == 0 && len(configPath) == 0 {
		return cli.NewExitError(errNoPath, 1)
	}
	var pass *string
	if len(configPath) != 0 {
		cfg, err := ReadWalletConfig(configPath)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		path = cfg.Path
		pass = &cfg.Password
	}
	wall, err := wallet.NewWallet(path)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	if err := wall.Save(); err != nil {
		return cli.NewExitError(err, 1)
	}

	if ctx.Bool("account") {
		if err := createAccount(wall, pass); err != nil {
			return cli.NewExitError(err, 1)
		}
		defer wall.Close()
	}

	fmtPrintWallet(ctx.App.Writer, wall)
	fmt.Fprintf(ctx.App.Writer, "wallet successfully created, file location is %s\n", wall.Path())
	return nil
}

func readAccountInfo() (string, string, error) {
	name, err := input.ReadLine("Enter the name of the account > ")
	if err != nil {
		return "", "", err
	}
	phrase, err := readNewPassword()
	if err != nil {
		return "", "", err
	}
	return name, phrase, nil
}

func readNewPassword() (string, error) {
	phrase, err := input.ReadPassword(EnterNewPasswordPrompt)
	if err != nil {
		return "", fmt.Errorf("Error reading password: %w", err)
	}
	phraseCheck, err := input.ReadPassword(ConfirmPasswordPrompt)
	if err != nil {
		return "", fmt.Errorf("Error reading password: %w", err)
	}

	if phrase != phraseCheck {
		return "", errPhraseMismatch
	}
	return phrase, nil
}

func createAccount(wall *wallet.Wallet, pass *string) error {
	var (
		name, phrase string
		err          error
	)
	if pass == nil {
		name, phrase, err = readAccountInfo()
		if err != nil {
			return err
		}
	} else {
		phrase = *pass
	}
	return wall.CreateAccount(name, phrase)
}

func openWallet(ctx *cli.Context, canUseWalletConfig bool) (*wallet.Wallet, *string, error) {
	path, pass, err := getWalletPathAndPass(ctx, canUseWalletConfig)
	if err != nil {
		return nil, nil, err
	}
	if path == "-" {
		return nil, nil, errNoStdin
	}
	w, err := wallet.NewWalletFromFile(path)
	if err != nil {
		return nil, nil, err
	}
	return w, pass, nil
}

func readWallet(ctx *cli.Context) (*wallet.Wallet, *string, error) {
	path, pass, err := getWalletPathAndPass(ctx, true)
	if err != nil {
		return nil, nil, err
	}
	if path == "-" {
		w := &wallet.Wallet{}
		if err := json.NewDecoder(os.Stdin).Decode(w); err != nil {
			return nil, nil, fmt.Errorf("js %w", err)
		}
		return w, nil, nil
	}
	w, err := wallet.NewWalletFromFile(path)
	if err != nil {
		return nil, nil, err
	}
	return w, pass, nil
}

// getWalletPathAndPass retrieves wallet path from context or from wallet configuration file.
// If wallet configuration file is specified, then account password is returned.
func getWalletPathAndPass(ctx *cli.Context, canUseWalletConfig bool) (string, *string, error) {
	path, configPath := ctx.String("wallet"), ctx.String("wallet-config")
	if !canUseWalletConfig && len(configPath) != 0 {
		return "", nil, errors.New("can't use wallet configuration file for this command")
	}
	if len(path) != 0 && len(configPath) != 0 {
		return "", nil, errConflictingWalletFlags
	}
	if len(path) == 0 && len(configPath) == 0 {
		return "", nil, errNoPath
	}
	var pass *string
	if len(configPath) != 0 {
		cfg, err := ReadWalletConfig(configPath)
		if err != nil {
			return "", nil, err
		}
		path = cfg.Path
		pass = &cfg.Password
	}
	return path, pass, nil
}

func ReadWalletConfig(configPath string) (*config.Wallet, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read wallet config: %w", err)
	}

	cfg := &config.Wallet{}

	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallet config YAML: %w", err)
	}
	return cfg, nil
}

func newAccountFromWIF(w io.Writer, wif string, scrypt keys.ScryptParams) (*wallet.Account, error) {
	// note: NEP2 strings always have length of 58 even though
	// base58 strings can have different lengths even if slice lengths are equal
	if len(wif) == 58 {
		pass, err := input.ReadPassword(EnterPasswordPrompt)
		if err != nil {
			return nil, fmt.Errorf("Error reading password: %w", err)
		}

		return wallet.NewAccountFromEncryptedWIF(wif, pass, scrypt)
	}

	acc, err := wallet.NewAccountFromWIF(wif)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "Provided WIF was unencrypted. Wallet can contain only encrypted keys.")
	name, pass, err := readAccountInfo()
	if err != nil {
		return nil, err
	}

	acc.Label = name
	if err := acc.Encrypt(pass, scrypt); err != nil {
		return nil, err
	}

	return acc, nil
}

func addAccountAndSave(w *wallet.Wallet, acc *wallet.Account) error {
	for i := range w.Accounts {
		if w.Accounts[i].Address == acc.Address {
			return fmt.Errorf("address '%s' is already in wallet", acc.Address)
		}
	}

	w.AddAccount(acc)
	return w.Save()
}

func fmtPrintWallet(w io.Writer, wall *wallet.Wallet) {
	b, _ := wall.JSON()
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, string(b))
	fmt.Fprintln(w, "")
}
