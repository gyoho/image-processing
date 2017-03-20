# Image Processing with Golang
A RESTful image processing web application server in GO!

## Proposed Architecture
Web server + SQS + App server (auto-scaling group) + MongoDB + S3
![architecure image](https://github.com/gyoho/image-processing/blob/master/architecture_diagram.jpg)

## Functionalities
* Receive an input image from a client, convert it to grayscale, compute the histogram and store both the original image and the histogram into a database
* Extract the histograms of the current week for a single user
* Extract the median histogram of the current day for all users
* Given a user id as input, returns n user ids who have the most similar histograms with respect to the input user id
