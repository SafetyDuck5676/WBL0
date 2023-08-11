// Function to handle the button click event
function handleClick() {
  var input = document.getElementById("inputField").value;
  var output = document.getElementById("response");

  fetch("/getOrder?data="+input, {
      method: "GET",
      headers: {
          "X-Requested-With": "XMLHttpRequest",
      },
  })
      .then(response => response.json())
      .then(response => {
          // Handle the response data
          var responseString = JSON.stringify(response);
          innerHTML = "<table>"
          innerHTML += "<tr><th>Order UID</th><th>entry</th><th>Customer</th><th>Track Number</th><th>Delivery Service</th></tr>"
          innerHTML += "<tr><td>"+response.order_uid+"</td><td>"+response.entry+"</td><td>"+response.customer_id+"</td><td>"+response.track_number+"</td><td>"+response.delivery_service+"</td></tr>"
          innerHTML += "</table>"
          output.innerHTML = innerHTML;
      })
      .catch(error => {
          // Handle any errors
          console.error(error);
      });
}

document.getElementById("submitButton").addEventListener("click", handleClick);