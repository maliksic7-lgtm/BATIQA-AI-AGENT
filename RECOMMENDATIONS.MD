Recommendations
GET
/api/recommendations

Optional:

/api/recommendations?category=restaurant
/api/recommendations?max_price=100000

Response:

{
  "items": [
    {
      "name": "Example Restaurant",
      "category": "restaurant",
      "price_min": 30000,
      "price_max": 90000,
      "distance_km": 1.2
    }
  ]
}