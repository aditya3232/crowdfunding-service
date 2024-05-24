package payment_gateway

import (
	"crowdfunding-service/internal/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/veritrans/go-midtrans"
)

type MidtransSnapGateway struct {
	Log    *logrus.Logger
	Config *viper.Viper
}

func NewMidtransSnapGateway(log *logrus.Logger, config *viper.Viper) *MidtransSnapGateway {
	return &MidtransSnapGateway{
		Log:    log,
		Config: config,
	}
}

// get payment url
func (m *MidtransSnapGateway) GetPaymentURL(paymentRequest *model.PaymentRequest) (string, error) {
	client := midtrans.NewClient()
	client.ServerKey = m.Config.GetString("midtrans.serverKey")
	client.ClientKey = m.Config.GetString("midtrans.clientKey")
	client.APIEnvType = midtrans.Sandbox // set to sandbox environment for testing

	snapGateway := midtrans.SnapGateway{
		Client: client,
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  paymentRequest.TransactionID,
			GrossAmt: int64(paymentRequest.Amount),
		},
		CustomerDetail: &midtrans.CustDetail{
			Email: paymentRequest.Email,
			FName: paymentRequest.FullName,
		},
	}

	snapToken, err := snapGateway.GetToken(snapReq)
	if err != nil {
		m.Log.Error("Error getting snap token:", err)
		return "", err
	}

	return snapToken.RedirectURL, nil
}
