package order

import (
	"beli-mang/internal/entity"
	"log"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Service interface {
	Create(orders []entity.Order, req CreateEstimationRequest) (resp CreateEstimationResponse, err error)
	CreateOrder(req CreateOrderRequest) (resp CreateOrderResponse, err error)
	Query(params QueryOrdersRequest) (items []QueryOrdersResponse, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Query(params QueryOrdersRequest) (items []QueryOrdersResponse, err error) {
	items, err = s.repo.Query(params)
	if err != nil {
		return items, err
	}

	if len(items) == 0 {
		return []QueryOrdersResponse{}, nil
	}

	return items, nil
}

func (s *service) CreateOrder(req CreateOrderRequest) (resp CreateOrderResponse, err error) {
	orderId, err := s.repo.CreateOrder(req)
	if err != nil {
		return resp, err
	}

	resp.OrderId = orderId

	return resp, nil
}

func (s *service) Create(orders []entity.Order, req CreateEstimationRequest) (resp CreateEstimationResponse, err error) {
	totalPrice, err := s.getTotalPrice(orders)
	if err != nil {
		return resp, err
	}

	// get starting order
	var startingMerchantId uuid.UUID
	for _, order := range req.Orders {
		if order.IsStartingPoint {
			uuid, err := uuid.Parse(order.MerchantId)
			if err != nil {
				return resp, fiber.NewError(fiber.StatusBadRequest, "invalid merchant starting point id")
			}

			startingMerchantId = uuid
		}
	}

	deliveryTime, err := s.getDeliveryTime(orders, startingMerchantId, req.UserLocation)
	log.Println("delivery time", deliveryTime)
	if err != nil {
		return resp, err
	}

	estimation := entity.Estimation{
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: deliveryTime,
	}
	estimationId, err := s.repo.Create(orders, estimation)
	if err != nil {
		return resp, err
	}

	resp.EstimatedDeliveryTimeInMin = deliveryTime
	resp.CalculatedEstimateId = estimationId
	resp.TotalPrice = totalPrice

	return resp, nil
}

func (s *service) getTotalPrice(orders []entity.Order) (totalPrice int, err error) {
	// get item uuids from req
	var itemUuids []uuid.UUID
	for _, order := range orders {
		for _, item := range order.OrderItems {
			itemUuids = append(itemUuids, item.ItemId)
		}
	}

	// get items form uuids
	items, err := s.repo.GetItems(itemUuids)
	if err != nil {
		return 0, err
	}

	// create product map for faster search
	var itemsMap = make(map[uuid.UUID]entity.Item)
	for _, item := range items {
		itemsMap[item.ID] = item
	}

	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			itemPrice := itemsMap[orderItem.ItemId].Price
			totalPrice += itemPrice * orderItem.Quantity
		}
	}

	return totalPrice, nil
}

type Location struct {
	Lat  float64
	Long float64
}

func (s *service) getDeliveryTime(orders []entity.Order, startingMerchantId uuid.UUID, userLocation UserLocation) (timeInMinutes int, err error) {
	var merchantIds []uuid.UUID
	for _, order := range orders {
		merchantIds = append(merchantIds, order.MerchantId)
	}

	// get merchants from uuids
	merchants, err := s.repo.GetMerchants(merchantIds)
	if err != nil {
		return 0, err
	}

	// create merchant map for faster search
	var merchantsMap = make(map[uuid.UUID]entity.Merchant)
	for _, merchant := range merchants {
		merchantsMap[merchant.ID] = merchant
	}

	// Starting point
	start := merchantsMap[startingMerchantId]

	// Initialize visited map
	visited := make(map[string]bool)
	visited[start.Name] = true

	// Route array
	route := []Location{}
	route = append(route, Location{start.Latitude, start.Longitude})

	// Find the route using Nearest Neighbor
	current := start
	for len(route) < len(merchants) {
		nearest := nearestNeighbor(current, merchants, visited)
		route = append(route, Location{
			nearest.Latitude,
			nearest.Longitude,
		})
		visited[nearest.Name] = true
		current = nearest
	}

	route = append(route, Location(userLocation))

	// Calculate the total distance of the route
	totalDistance := getTotalDistance(route)
	log.Println("totalDistance", totalDistance)

	// Calculate the time in minutes
	deliverySpeed := 40.0 //Km per hour
	estimatedTime := totalDistance / deliverySpeed
	estimatedTime = estimatedTime * 60

	return int(estimatedTime), nil
}

// Haversine formula to calculate distance between two points
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of Earth in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// Find the nearest unvisited merchant
func nearestNeighbor(current entity.Merchant, merchants []entity.Merchant, visited map[string]bool) entity.Merchant {
	minDist := math.MaxFloat64
	var nearest entity.Merchant
	for _, merchant := range merchants {
		if !visited[merchant.Name] {
			dist := haversine(current.Latitude, current.Longitude,
				merchant.Latitude, merchant.Longitude)
			if dist < minDist {
				minDist = dist
				nearest = merchant
			}
		}
	}
	return nearest
}

// Calculate the total distance of the route
func getTotalDistance(route []Location) float64 {
	totalDist := 0.0
	for i := 0; i < len(route)-1; i++ {
		totalDist += haversine(route[i].Lat, route[i].Long, route[i+1].Lat, route[i+1].Long)
	}
	return totalDist
}
