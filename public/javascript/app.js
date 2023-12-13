document.addEventListener('DOMContentLoaded', function () {
  const cells = document.querySelectorAll('[data-habit-id][data-day-index]');

  cells.forEach(cell => {
      cell.addEventListener('click', function () {
          const habitId = this.getAttribute('data-habit-id');
          const dayIndex = this.getAttribute('data-day-index');
          console.log('dayIndex:', dayIndex);


          // Check the current background color
          const isCompleted = this.style.backgroundColor === 'lightgreen';

          // Send an AJAX request to update Repeat_Count
          fetch(`/habit/updateRepeatCount/${habitId}`, {
              method: 'PUT',
              headers: {
                  'Content-Type': 'application/json',
              },
              body: JSON.stringify({ dayIndex: parseInt(dayIndex), completed: !isCompleted }),
          })
          .then(response => response.json())
          .then(data => {
              // Update the cell content
              this.textContent = data.message;

              // Toggle the cell color
              this.style.backgroundColor = isCompleted ? '' : 'lightgreen';
          })
          .catch(error => console.error('Error:', error));
      });
  });
});
