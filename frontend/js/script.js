import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import Chart from 'chart.js/auto';

/* async function fetchDataAndRenderChart() {
  try {
    const response = await fetch("/list");
    const labelText = await response.text(); // Await the response text

    console.log(labelText)

    const ctx = document.getElementById('myChart');
    if (!ctx) {
      console.error("Canvas element with ID 'myChart' not found.");
      return;
    }

    new Chart(ctx, {
      type: 'bar',
      data: {
        labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
        datasets: [{
          label: labelText, // Corrected
          data: [12, 19, 3, 5, 2, 3],
          backgroundColor: 'black', // Bars will be black
          borderColor: 'black',     // Optional: Black border
          borderWidth: 1
        }]
      },
      options: {
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              color: 'black' // Y-axis font color
            }
          },
          x: {
            ticks: {
              color: 'black' // X-axis font color
            }
          }
        },
        plugins: {
          legend: {
            labels: {
              color: 'black' // Legend font color
            }
          }
        }
      }
    });
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

// Call the function to fetch data and render the chart
fetchDataAndRenderChart();*/

const ctx = document.getElementById('myChart');
if (!ctx) {
  console.error("Canvas element with ID 'myChart' not found.");
}

new Chart(ctx, {
  type: 'bar',
  data: {
    labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
    datasets: [{
      label: "Users",
      data: [12, 19, 3, 5, 2, 3],
      backgroundColor: 'black', // Bars will be black
      borderColor: 'black',     // Optional: Black border
      borderWidth: 1
    }]
  },
  options: {
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          color: 'black' // Y-axis font color
        }
      },
      x: {
        ticks: {
          color: 'black' // X-axis font color
        }
      }
    },
    plugins: {
      legend: {
        labels: {
          color: 'black' // Legend font color
        }
      }
    }
  }
});

  