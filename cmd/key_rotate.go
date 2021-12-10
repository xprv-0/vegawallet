package cmd

import (
	"fmt"
	"io"

	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	rotateKeyLong = cli.LongDesc(`
		Build a signed key rotation transaction as a Base64 encoded string.
		Choose a public key to rotate to and target block height.

		The generated transaction can be sent using the command: "tx send".
	`)

	rotateKeyExample = cli.Examples(`
		# Build signed transaction for rotating to new key public key
		vegawallet key rotate --wallet WALLET --tx-height TX_HEIGHT --target-height TARGET_HEIGHT --pubkey PUBLIC_KEY --current-pubkey CURRENT_PUBLIC_KEY
	`)
)

type RotateKeyHandler func(*wallet.RotateKeyRequest) (*wallet.RotateKeyResponse, error)

func NewCmdRotateKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.RotateKeyRequest) (*wallet.RotateKeyResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.RotateKey(s, req)
	}

	return BuildCmdRotateKey(w, h, rf)
}

func BuildCmdRotateKey(w io.Writer, handler RotateKeyHandler, rf *RootFlags) *cobra.Command {
	f := RotateKeyFlags{}

	cmd := &cobra.Command{
		Use:     "rotate",
		Short:   "Build a signed key rotation transaction",
		Long:    rotateKeyLong,
		Example: rotateKeyExample,
		RunE: func(_ *cobra.Command, args []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}

			resp, err := handler(req)
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintRotateKeyResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"Wallet holding the master key and new public key",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)
	cmd.Flags().StringVar(&f.NewPublicKey,
		"new-pubkey",
		"",
		"A public key to rotate to. Should be generated by wallet's 'generate' command.",
	)
	cmd.Flags().StringVar(&f.CurrentPubKey,
		"current-pubkey",
		"",
		"A public key to rotate from. Should be currently used public key.",
	)
	cmd.Flags().Uint64Var(&f.TxBlockHeight,
		"tx-height",
		0,
		"It should be close to the current block height when the transaction is applied, with a threshold of ~ - 150 blocks.",
	)
	cmd.Flags().Uint64Var(&f.TargetBlockHeight,
		"target-height",
		0,
		"Height of block where the public key change will take effect",
	)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type RotateKeyFlags struct {
	Wallet            string
	PassphraseFile    string
	NewPublicKey      string
	CurrentPubKey     string
	TxBlockHeight     uint64
	TargetBlockHeight uint64
}

func (f *RotateKeyFlags) Validate() (*wallet.RotateKeyRequest, error) {
	req := &wallet.RotateKeyRequest{}

	if f.NewPublicKey == "" {
		return nil, flags.FlagMustBeSpecifiedError("new-pubkey")
	}
	req.NewPublicKey = f.NewPublicKey

	if f.CurrentPubKey == "" {
		return nil, flags.FlagMustBeSpecifiedError("current-pubkey")
	}
	req.CurrentPublicKey = f.CurrentPubKey

	if f.TargetBlockHeight == 0 {
		return nil, flags.FlagMustBeSpecifiedError("target-height")
	}
	req.TargetBlockHeight = f.TargetBlockHeight

	if f.TxBlockHeight == 0 {
		return nil, flags.FlagMustBeSpecifiedError("tx-height")
	}
	req.TxBlockHeight = f.TxBlockHeight

	if req.TargetBlockHeight <= req.TxBlockHeight {
		return nil, flags.FlagRequireLessThanFlagError("tx-height", "target-height")
	}

	if len(f.Wallet) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("wallet")
	}
	req.Wallet = f.Wallet

	passphrase, err := flags.GetPassphrase(f.PassphraseFile)
	if err != nil {
		return nil, err
	}
	req.Passphrase = passphrase

	return req, nil
}

func PrintRotateKeyResponse(w io.Writer, req *wallet.RotateKeyResponse) {
	p := printer.NewInteractivePrinter(w)
	p.CheckMark().SuccessText("Key rotation succeeded").NextSection()
	p.Text("Base64 encoded transaction:").NextLine()
	p.Text(req.Base64Transaction).NextLine()
	p.Text("Master public key used:").NextLine()
	p.Text(req.MasterPublicKey).NextLine()
}
