package cmd

import (
	"bytes"
	"context"
	"github.com/k-yomo/pubsub_cli/util"
	"github.com/spf13/cobra"
	"testing"
	"time"
)

func Test_registerPush(t *testing.T) {
	pubsubClient, err := util.NewTestPubSubClient(t)
	if err != nil {
		t.Fatal(err)
	}
	clear := setTestRootVariables(t)
	defer clear()

	type args struct {
		in0          *cobra.Command
		pubsubClient *util.PubSubClient
		args         []string
	}
	tests := []struct {
		name               string
		mockSubscriptionID string
		args               args
		check              func()
		wantErr            bool
	}{
		{
			name:               "push subscription is registered successfully",
			mockSubscriptionID: "test",
			args:               args{pubsubClient: pubsubClient, args: []string{"test_topic", "http://localhost:9000"}},
			check: func() {
				sub := pubsubClient.Subscription("test")
				subConfig, err := sub.Config(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				topic := "test_topic"
				// check if topic is collect
				if subConfig.Topic.ID() != topic {
					t.Errorf("registerPush() got topic = %v, want %v", subConfig.Topic.String(), topic)
				}
				// check if endpoint is collect
				if subConfig.PushConfig.Endpoint != "http://localhost:9000" {
					t.Errorf("registerPush() got endpoint = %v, want %v", subConfig.PushConfig.Endpoint, "http://localhost:9000")
				}
				// check if expirationPolicy is set to 24 hours
				if subConfig.ExpirationPolicy != 24*time.Hour {
					t.Errorf("registerPush() got expirationPolicy = %v, want %v", subConfig.ExpirationPolicy, 24*time.Hour)
				}
				sub.Delete(context.Background())
			},
		},
		{
			name:    "push subscription with invalid topic name causes error",
			args:    args{pubsubClient: pubsubClient, args: []string{"1", "http://localhost:9000"}},
			check:   func() {},
			wantErr: true,
		},
		{
			name:    "push subscription with invalid endpoint causes error",
			args:    args{pubsubClient: pubsubClient, args: []string{"test_topic", "invalid"}},
			check:   func() {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clear := util.SetMockUUID(t, tt.mockSubscriptionID)
			defer clear()

			out := &bytes.Buffer{}
			err := newRegisterPushCmd(out).RunE(tt.args.in0, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerPush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.check()
		})
	}
}
