package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("../../html"),
	jet.InDevelopmentMode(),
)

var wsChan = make(chan WSJSONRequest)
var clients = make(map[WSConn]string)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

type WSConn struct {
	*websocket.Conn
}

type WSJSONRequest struct {
	Action   string `json:"action"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Conn     WSConn `json:"-"`
}

type WSJSONResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

func WSEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	var res WSJSONResponse
	res.Message = `<em><small>Connected to server</small></em>`

	conn := WSConn{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(res)
	if err != nil {
		log.Println(err)
	}

	go listenForWS(&conn)
}

func listenForWS(conn *WSConn) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var request WSJSONRequest

	for {
		err := conn.ReadJSON(&request)
		if err == nil {
			request.Conn = *conn
			wsChan <- request
		}
	}
}

func ListenToWSChannel() {
	var response WSJSONResponse

	for {
		request := <-wsChan

		switch request.Action {
		case "username":
			clients[request.Conn] = request.Username
			users := getUserList()
			response.Action = "user_list"
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "left":
			delete(clients, request.Conn)
			users := getUserList()
			response.Action = "user_list"
			response.ConnectedUsers = users
			broadcastToAll(response)
		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", request.Username, request.Message)
			broadcastToAll(response)
		}
	}
}

func getUserList() []string {
	var userList []string
	for _, user := range clients {
		if user != "" {
			userList = append(userList, user)
		}
	}
	sort.Strings(userList)
	return userList
}

func broadcastToAll(response WSJSONResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
