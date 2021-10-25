package cmd

import (
	"fmt"
	"os"
	"strings"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"github.com/spf13/cobra"

	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
)

var (
	keyGenerateArgs struct {
		wallet         string
		passphraseFile string
		metas          string
	}

	keyGenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a new key pair for a wallet",
		Long:  "Generate a new key pair for a wallet, this will implicitly generate a new wallet if none exist for the given wallet",
		RunE:  runKeyGenerate,
	}
)

func init() {
	keyCmd.AddCommand(keyGenerateCmd)
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.metas, "meta", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
	_ = keyGenerateCmd.MarkFlagRequired("wallet")
}

func runKeyGenerate(_ *cobra.Command, _ []string) error {
	p := printer.NewHumanPrinter()

	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	walletExists := handler.WalletExists(keyGenerateArgs.wallet)

	passphrase, err := getPassphrase(keyGenerateArgs.passphraseFile, !walletExists)
	if err != nil {
		return err
	}

	metas, err := parseMeta(keyGenerateArgs.metas)
	if err != nil {
		return err
	}

	var mnemonic string
	if !walletExists {
		if rootArgs.output == "human" {
			p.BangMark().Text("Wallet ").Bold(keyGenerateArgs.wallet).Text(" does not exist yet").NextLine()
		}

		mnemonic, err = handler.CreateWallet(keyGenerateArgs.wallet, passphrase)
		if err != nil {
			return fmt.Errorf("couldn't create wallet: %w", err)
		}

		if rootArgs.output == "human" {
			p.CheckMark().Text("Wallet ").Bold(keyGenerateArgs.wallet).Text(" has been created at: ").SuccessText(store.GetWalletPath(keyGenerateArgs.wallet)).NextLine()
		}
	}

	keyPair, err := handler.GenerateKeyPair(keyGenerateArgs.wallet, passphrase, metas)
	if err != nil {
		return fmt.Errorf("could not generate a key pair: %w", err)
	}

	if rootArgs.output == "human" {
		printHuman(p, mnemonic, keyPair, store.GetWalletPath(keyGenerateArgs.wallet))
	} else if rootArgs.output == "json" {
		return printKeyGenerateJSON(mnemonic, keyPair, store.GetWalletPath(keyGenerateArgs.wallet))
	} else {
		return NewUnsupportedCommandOutputError(rootArgs.output)
	}
	return nil
}

func printHuman(p *printer.HumanPrinter, mnemonic string, keyPair wallet.KeyPair, walletPath string) {
	p.CheckMark().Text("Key pair has been generated for wallet ").Bold(keyGenerateArgs.wallet).Text(" at: ").SuccessText(walletPath).NextLine()
	p.CheckMark().SuccessText("Generating a key pair succeeded").NextSection()
	if len(mnemonic) != 0 {
		p.Text("Wallet mnemonic:").NextLine().WarningText(mnemonic).NextLine()
	}
	p.Text("Public key:").NextLine().WarningText(keyPair.PublicKey()).NextLine()
	p.Text("Metadata:").NextLine()
	printMeta(p, keyPair.Meta())
	p.NextLine()

	p.RedArrow().DangerText("Important").NextLine()
	if len(mnemonic) != 0 {
		p.DangerText("1. ").Text("Write down the mnemonic and store it somewhere safe and secure, now, as it will ").Underline("not").Text(" be displayed ever again!").NextLine()
		p.DangerText("2. ").Text("Do not share the mnemonic nor the private key.").NextSection()
	} else {
		p.Text("Do not share the mnemonic nor the private key.").NextSection()
	}

	p.BlueArrow().InfoText("Run the service").NextLine()
	p.Text("Once you have a key pair generated, you can run the service with the following command:").NextSection()
	p.Code(fmt.Sprintf("%s service run --network <NETWORK_TO_CONNECT_TO>", os.Args[0])).NextSection()
	p.Text("If you want to open up a local version of Vega Console alongside the service, use the following command:").NextSection()
	p.Code(fmt.Sprintf("%s service run --network <NETWORK_TO_CONNECT_TO> --console-proxy", os.Args[0])).NextSection()
	p.Text("To terminate the process, hit ").Bold("ctrl+c").NextSection()
	p.Text("For more information, use ").Bold("--help").Text(" flag.").NextLine()
}

type keyGenerateJSON struct {
	Wallet keyGenerateWalletJSON `json:"wallet"`
	Key    keyGenerateKeyJSON    `json:"key"`
}

type keyGenerateWalletJSON struct {
	FilePath string `json:"filePath"`
	Mnemonic string `json:"mnemonic,omitempty"`
}

type keyGenerateKeyJSON struct {
	KeyPair   keyGenerateKeyPairJSON   `json:"keyPair"`
	Algorithm keyGenerateAlgorithmJSON `json:"algorithm"`
	Meta      []wallet.Meta            `json:"meta"`
}

type keyGenerateKeyPairJSON struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type keyGenerateAlgorithmJSON struct {
	Name    string `json:"name"`
	Version uint32 `json:"version"`
}

func printKeyGenerateJSON(mnemonic string, keyPair wallet.KeyPair, walletPath string) error {
	result := keyGenerateJSON{
		Wallet: keyGenerateWalletJSON{
			FilePath: walletPath,
			Mnemonic: mnemonic,
		},
		Key: keyGenerateKeyJSON{
			KeyPair: keyGenerateKeyPairJSON{
				PrivateKey: keyPair.PrivateKey(),
				PublicKey:  keyPair.PublicKey(),
			},
			Algorithm: keyGenerateAlgorithmJSON{
				Name:    keyPair.AlgorithmName(),
				Version: keyPair.AlgorithmVersion(),
			},
			Meta: keyPair.Meta(),
		},
	}
	return vgjson.Print(result)
}

func parseMeta(metaStr string) ([]wallet.Meta, error) {
	if len(metaStr) == 0 {
		return nil, nil
	}

	rawMetas := strings.Split(metaStr, ";")

	metas := make([]wallet.Meta, 0, len(rawMetas))
	for _, v := range rawMetas {
		rawMeta := strings.Split(v, ":")
		if len(rawMeta) != 2 { //nolint:gomnd
			return nil, ErrInvalidMetadataFormat
		}
		metas = append(metas, wallet.Meta{Key: rawMeta[0], Value: rawMeta[1]})
	}

	return metas, nil
}
