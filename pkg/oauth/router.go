package oauth

type Container struct {
	clients map[Provider]Client
}

func NewContainer(clients ...Client) *Container {
	container := &Container{
		clients: make(map[Provider]Client),
	}
	for _, client := range clients {
		container.clients[client.GetProvider()] = client
	}
	return container
}

func (container *Container) Use(provider Provider) Client {
	return container.clients[provider]
}
