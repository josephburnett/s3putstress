package main

import (
        "fmt"
        "math/rand"
        "time"

        "gopkg.in/amz.v1/aws"
        "gopkg.in/amz.v1/s3"
)

func main() {

        // We want the put the same sequence of objects so we don't accumulate storage
        rand.Seed(0)

        // Limit determines how fast objects are put.
        limit := make(chan bool)
        goHome := make(chan bool)

        // 100 partiers :)
        for i := 0; i < 100; i++ {
                go putParty(limit);
        }

        // 1 party pooper :(
        go partyPooper(limit, goHome)
        <- goHome
}

func putParty(limit chan bool) {

        // Put to the max!
        for {
                <-limit
                put()
        }
}

func partyPooper(limit, goHome chan bool) {

        // Don't party too fast!  Don't party too long!
        start := time.Now()
        for {
                elapsed := time.Since(start)
                if (elapsed > 1000000000 * 60 * 10 /* 10 minutes */) {
                        fmt.Printf("\n")
                        goHome <- true
                }

                limit <- true
                time.Sleep(10 * time.Millisecond) /* 100 TPS */
        }
}

func put() {

        // The AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables are used.
        auth, err := aws.EnvAuth()
        if err != nil {
                panic(err.Error())
        }

        // Open Bucket
        s := s3.New(auth, aws.USEast)
        bucket := s.Bucket("joburnet-lambda-sources")

        data := []byte("bang!")
        key := randSeq(32)

        err = bucket.Put(key, data, "content-type", s3.Private)
        if err != nil {
                fmt.Printf("%v", err)
        } else {
                fmt.Printf(".")
        }
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
        b := make([]rune, n)
        for i := range b {
                b[i] = letters[rand.Intn(len(letters))]
        }
        return string(b)
}
