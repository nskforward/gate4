package handler

/*
func CreateCert(client *api.Client) Handler {
	return func(ctx context.Context, args []string) error {
		scanner := console.NewScanner()
		defer scanner.Close()

		commonName, err := scanner.Scan(ctx, "common name", "", nil)
		if err != nil {
			return err
		}

		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = pem.Encode(&buf, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
		if err != nil {
			return err
		}

		privateKey := buf.String()

		cert, err := client.CreateCert(ctx, commonName, privateKey)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("success: copy cert below")

		fmt.Println()
		fmt.Println(privateKey)
		fmt.Println()
		fmt.Println(cert)
		fmt.Println()

		return nil
	}
}
*/
