# IRCTC-TrainSearch-GoLang-API

#Here is the step-by-step textual flowchart that describes the product functionality and system 
#design of the IRCTC Train Search page (https://www.irctc.co.in/nget/train-search): 

IRCTC Train Search – Step-by-Step Functionality & System Flow 
1. User opens the IRCTC Train Search page. → URL: https://www.irctc.co.in/nget/train-search 
2. The frontend (Angular Single Page App) loads in the user's browser. 
3. User enters input: 
  ○ From Station (source) 
  ○ To Station (destination) 
  ○ Journey Date 
  ○ Class (e.g., Sleeper, 3AC, 2AC) 
  ○ Quota (e.g., General, Tatkal) 
4. User clicks the “Search” button. 
5. Frontend creates an API request to IRCTC backend service (Train Availability API). → 
Request includes: 
  ○ fromStationCode 
  ○ toStationCode 
  ○ journeyDate 
  ○ classCode 
  ○ quotaCode 
6. Backend validates inputs and checks: 
  ○ Train schedule between source and destination 
  ○ Class-wise seat availability 
  ○ Fare and travel time 
  ○ Booking rules (quota, day restrictions) 
7. Backend returns a JSON response containing: 
  ○ List of trains 
  ○ Train numbers & names 
  ○ Departure & arrival times 
  ○ Running days 
  ○ Travel duration 
  ○ Class-wise availability 
  ○ Fare per class 
  ○ Bookability status 
8. Frontend displays search results with filtering options: 
  ○ Sort by departure time, duration, fare, availability 
  ○ Filter by class, quota, train type, date 
9. User selects a train and clicks “Book Now”. 
10. System checks if user is logged in: 
  ○ If not → redirect to login page 
  ○ If logged in → redirect to passenger details page 
11. Booking flow continues: 
  ○ Enter passenger details 
  ○ Choose payment method 
  ○ Confirm booking 

 Key Components Involved: 
● UI: Angular Frontend (SPA) 
● API Gateway / Load Balancer 
● Backend Microservices: 
  ○ Train Schedule Service
  ○ Seat Availability Service 
  ○ Fare Calculation Service 
  ○ Booking Service 
● Database: 
  ○ Train Timetable DB 
  ○ Reservation DB 
  ○ Fare Master Tables
