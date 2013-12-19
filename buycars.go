package main

import (
  "net/http"
  "fmt"
  "github.com/garyburd/redigo/redis"
  "log"
  "errors"
)

type Item struct {
  Title string
  Description string
  Image []byte
}

func (i *Item) String() string {
  return i.Title
}

func loadItem(key string) (*Item, error) {
  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    return nil, err
  }

  reply, err := redis.Values(c.Do("HGETALL", key))
  if err != nil {
    return nil, err
  }
  if len(reply) == 0 {
    return nil, errors.New("No Item Found")
  }

  item := &Item{}
  err = redis.ScanStruct(reply, item)
  if err != nil {
    return nil, err
  }
  return item, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  item := r.URL.Path[len("/view/"):]
  i, err := loadItem(item)
  if err != nil {
    log.Print(err)
    http.NotFound(w, r)
    return
  }
  fmt.Fprintf(w, "%s", i)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.ListenAndServe(":8080", nil)
}
