package main

import (
        "fmt"
        "math/rand"
        "time"

        "gopkg.in/amz.v1/aws"
        "gopkg.in/amz.v1/s3"
)

var bucket *s3.Bucket

func main() {

        // The AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables are used.
        auth, err := aws.EnvAuth()
        if err != nil {
                panic(err.Error())
        }

        // Open Bucket
        s := s3.New(auth, aws.USEast)
        bucket = s.Bucket("joburnet-lambda-sources")

        rand.Seed(0)

        limit := make(chan bool)
        done := make(chan bool)

        for i := 0; i < 50; i++ {
                go putParty(limit);
        }

        go partyPooper(limit, done)
        <- done
}

func putParty(limit chan bool) {
        for {
                <-limit
                put()
        }
}

func partyPooper(limit, done chan bool) {

        for i := 0; i < 100; i++ {
                limit <- true
                time.Sleep(100 * time.Millisecond)
        }

        fmt.Printf("\n")
        done <- true
}

func put() {

        data := []byte("bang!")
        key := randSeq(32)

        err := bucket.Put(key, data, "content-type", s3.Private)
        if err != nil {
                panic(err)
        }

        fmt.Printf(".")
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
        b := make([]rune, n)
        for i := range b {
                b[i] = letters[rand.Intn(len(letters))]
        }
        return string(b)
}
