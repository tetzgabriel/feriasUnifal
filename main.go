package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func main() {
	lambda.Start(executeBot)
}

func executeBot() {
	log.Printf("Starting FeriasUnifalBot...")

	log.Printf("Getting Credentials from environment...")
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN_FERIAS"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET_FERIAS"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY_FERIAS"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET_FERIAS"),
	}

	log.Printf("Getting Twitter client...\n")
	client, err := getClient(&creds)
	if err != nil {
		log.Fatal("--- Error getting Twitter Client, shutting down app :( --- ")
		log.Println(err)
	}

	phraseToTweet := getRandomPhrase()

	tweet(client, phraseToTweet)
}

func getRandomPhrase() string {
	phrases := [10]string{
		"Não pare de estudar, pelo menos você vai pegar uma DP sendo inteligente ;)",
		"Sem lutas não há derrotas!",
		"A faculdade é um grande lençol de elástico, quando você ajeita de um lado, ela solta de outro",
		"A faculdade é um conto de falhas",
		"Faculdade é igual Uno, você resolve um problema e depois vem +4",
		"Novos dias, novos erros",
		"Estamos todos no mesmo barco, ele chama Titanic",
		"Daqui pra frente é só pra trás",
		"Bota a cara no sol mona",
		"Já acabou Jéssica? Ainda não",
	}

	num := getRandomInt()

	return phrases[num]
}

func getRandomInt() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 9

	return rand.Intn(max-min+1) + min
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's ACCOUNT: %s\n", user.Name)
	return client, nil
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func getDaysUntilVacation() int {
	vacationDay := Date(2022, 4, 16)
	today := time.Now()

	days := vacationDay.Sub(today).Hours() / 24

	log.Printf("Days left: %d\n", int(days))
	return int(days)
}

func tweet(client *twitter.Client, phrase string) {

	daysUntilVacation := getDaysUntilVacation()
	tweetText := ""

	if daysUntilVacation == 0 {
		tweetText = "Começooooooooooo\n Todo mundo livreeeee\n Boas férias galeraaaaaaa!"
	} else if daysUntilVacation == 1 {
		tweetText = "Só amanhã, só amanhã mesmo\n As férias da 08 começam amanhã!\n Parabéns aos sobreviventes ;)"
	} else {
		tweetText = fmt.Sprintf("%s\nFaltam só %d dias para as férias da 08!", phrase, daysUntilVacation)
	}

	tweet, _, err := client.Statuses.Update(tweetText, nil)
	if err != nil {
		log.Printf("Error tweeting!")
		log.Fatal(err)
	}

	log.Printf("Success tweeting!")
	log.Printf("%+v\n", tweet)
}
