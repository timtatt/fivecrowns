package main

import (
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"net/http"

	"github.com/timtatt/fivecrowns/bots"
	"github.com/timtatt/fivecrowns/bots/bigbrainbot"
	"github.com/timtatt/fivecrowns/bots/grugbot"
	"github.com/timtatt/fivecrowns/bots/smoothbrainbot"
)

func main() {

	port := *flag.String("port", "3000", "specify the port of the http server")

	slog.Info("starting web server")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("pong"))
	})

	fs := http.FileServer(http.Dir("./arena"))
	mux.Handle("/arena/", http.StripPrefix("/arena/", fs))

	configureBots(mux)

	slog.Info("listening on port " + port)
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		log.Fatal(err)
	}

}

func configureBots(mux *http.ServeMux) {

	b := map[string]bots.Bot{
		"smoothbrainbot": smoothbrainbot.NewSmoothBrainBot(),
		"grugbot":        grugbot.NewGrugBot(),
		"bigbrainbot":    bigbrainbot.NewBigBrainBot(),
	}

	for botName, bot := range b {
		slog.Info("registering endpoint", "bot", botName)
		mux.HandleFunc("POST /bots/"+botName, func(res http.ResponseWriter, req *http.Request) {
			slog.Info("recieved request", "bot", botName)

			defer req.Body.Close()

			var botReq bots.BotRequest
			err := json.NewDecoder(req.Body).Decode(&botReq)

			if err != nil {
				slog.Error("unable to unmarshal request", "err", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			var botRes interface{}
			slog.Info("received request", "action", botReq.Action, "bot", botName, "req", botReq)
			switch botReq.Action {
			case bots.ActionScore:
				botRes, err = bot.Score(botReq)
			case bots.ActionDiscard:
				botRes, err = bot.Discard(botReq)
			case bots.ActionDraw:
				botRes, err = bot.Draw(botReq)
			}
			slog.Info("calculated response", "action", botReq.Action, "bot", botName, "res", botRes)

			if err != nil {
				slog.Error("failed to get bot response", "err", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(res).Encode(botRes)

			if err != nil {
				slog.Error("failed to encode bot response", "err", err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

		})
	}

}
