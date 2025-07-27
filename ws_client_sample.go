package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/akhilesh-saipangallu/stream_benchmarking/config"
	"github.com/akhilesh-saipangallu/stream_benchmarking/datatype"
	"github.com/akhilesh-saipangallu/stream_benchmarking/metrics"
	"github.com/akhilesh-saipangallu/stream_benchmarking/utils"
	"github.com/akhilesh-saipangallu/stream_benchmarking/ws"
	"github.com/gorilla/websocket"
)

func main() {
	config.Init()
	cfg, err := config.Get()
	if err != nil {
		log.Println("main: %s", err)
		return
	}

	url := cfg.WsServer.WSEndpoint

	rLogger := utils.NewLogOnceEveryN(100)

	log.Printf("connecting to %s", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer conn.Close()

	// Example: Send a subscription message
	subMsg := ws.SubscriptionMessage{
		Action: "subscribe",
		Tickers: []string{
			"bn_0",
			"bn_1",
			"bn_2",
			"bn_3",
			"bn_4",
			"bn_5",
			"bn_6",
			"bn_7",
			"bn_8",
			"bn_9",
			"bn_10",
			"bn_11",
			"bn_12",
			"bn_13",
			"bn_14",
			"bn_15",
			"bn_16",
			"bn_17",
			"bn_18",
			"bn_19",
			"bn_20",
			"bn_21",
			"bn_22",
			"bn_23",
			"bn_24",
			"bn_25",
			"bn_26",
			"bn_27",
			"bn_28",
			"bn_29",
			"bn_30",
			"bn_31",
			"bn_32",
			"bn_33",
			"bn_34",
			"bn_35",
			"bn_36",
			"bn_37",
			"bn_38",
			"bn_39",
			"bn_40",
			"bn_41",
			"bn_42",
			"bn_43",
			"bn_44",
			"bn_45",
			"bn_46",
			"bn_47",
			"bn_48",
			"bn_49",
			"bn_50",
			"bn_51",
			"bn_52",
			"bn_53",
			"bn_54",
			"bn_55",
			"bn_56",
			"bn_57",
			"bn_58",
			"bn_59",
			"bn_60",
			"bn_61",
			"bn_62",
			"bn_63",
			"bn_64",
			"bn_65",
			"bn_66",
			"bn_67",
			"bn_68",
			"bn_69",
			"bn_70",
			"bn_71",
			"bn_72",
			"bn_73",
			"bn_74",
			"bn_75",
			"bn_76",
			"bn_77",
			"bn_78",
			"bn_79",
			"bn_80",
			"bn_81",
			"bn_82",
			"bn_83",
			"bn_84",
			"bn_85",
			"bn_86",
			"bn_87",
			"bn_88",
			"bn_89",
			"bn_90",
			"bn_91",
			"bn_92",
			"bn_93",
			"bn_94",
			"bn_95",
			"bn_96",
			"bn_97",
			"bn_98",
			"bn_99",
		},
	}
	err = conn.WriteJSON(subMsg)
	if err != nil {
		log.Println("write error:", err)
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for {
			now := time.Now().UnixMicro()

			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			pm, err := datatype.ParsePriceMessageBytes(msg)
			if err != nil {
				return
			}
			latency := float64(now - pm.Timestamp) // latency in µs
			metrics.EndToEndLatency.Observe(latency)

			// log.Printf("received: %s", msg)
			rLogger.Log("received: %s", msg)
		}
	}()

	// Keep the client running until interrupted
	for {
		select {
		case <-interrupt:
			log.Println("interrupt received; closing connection...")

			// Optional: Send close message
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("close error:", err)
				return
			}
			time.Sleep(time.Second)
			return
		}
	}
}
