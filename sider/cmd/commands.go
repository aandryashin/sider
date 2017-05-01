package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aandryashin/sider/siderd/client"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	ttl time.Duration
)

func init() {
	setCmd.Flags().DurationVarP(&ttl, "ttl", "", 0, "key expiration timeout")

}

var (
	keysCmd = &cobra.Command{
		Use:   "keys",
		Short: "List of stored keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("wrong args number")
			}
			cl := &client.Client{siderURL}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			keys, err := cl.Keys(ctx)
			if err != nil {
				return fmt.Errorf("client: %v", err)
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "    ")
			err = encoder.Encode(keys)
			if err != nil {
				return fmt.Errorf("output keys: %v", err)
			}
			return nil
		},
	}
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get value by the key",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("missing key arg")
			}
			key := args[0]
			cl := &client.Client{siderURL}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			v, err := cl.Get(ctx, key)
			if err != nil {
				return fmt.Errorf("client: %v", err)
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "    ")
			err = encoder.Encode(v)
			if err != nil {
				return fmt.Errorf("output value: %v", err)
			}
			return nil
		},
	}
	setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set value for the key",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return fmt.Errorf("missing key and value args")
			case 1:
				return fmt.Errorf("missing value arg")
			default:
			}
			key, value := args[0], args[1]
			cl := &client.Client{siderURL}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			err := cl.Set(ctx, key, bytes.NewReader([]byte(value)), ttl)
			if err != nil {
				return fmt.Errorf("client: %v", err)
			}
			return nil
		},
	}
	delCmd = &cobra.Command{
		Use:   "del",
		Short: "Delete key and value",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("missing key arg")
			}
			key := args[0]
			cl := &client.Client{siderURL}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			err := cl.Del(ctx, key)
			if err != nil {
				return fmt.Errorf("client: %v", err)
			}
			return nil
		},
	}
)
