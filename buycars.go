package main

import (
  "net/http"
  "fmt"
  "github.com/garyburd/redigo/redis"
  "log"
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

  item := &Item{}
  err = redis.ScanStruct(reply, item)
  if err != nil {
    return nil, err
  }
  //TODO handle empty return value

  log.Printf("Item found: %v", item)
  return item, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  item := r.URL.Path[len("/view/"):]
  i, err := loadItem(item)
  if err != nil {
    log.Print(err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  fmt.Fprintf(w, "%s", i)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.ListenAndServe(":8080", nil)
}
