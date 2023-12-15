document.addEventListener('DOMContentLoaded', function () {
    const habitTable = document.getElementById('habitTable');

    habitTable.addEventListener('click', function (event) {
        const target = event.target;

        // Check if the clicked cell is in the habit table body
        if (target.tagName === 'TD' && target.parentNode.classList.contains('habit-row')) {
            const habitId = target.parentNode.dataset.habitId;
            const dayOfWeek = target.dataset.dayOfWeek;

            // Define the day index based on your requirements
            let dayIndex;
            switch (dayOfWeek) {
                case 'Monday':
                    dayIndex = 0;
                    break;
                case 'Tuesday':
                    dayIndex = 1;
                    break;
                case 'Wednesday':
                    dayIndex = 2;
                    break;
                case 'Thursday':
                    dayIndex = 3;
                    break;
                case 'Friday':
                    dayIndex = 4;
                    break;
                case 'Saturday':
                    dayIndex = 5;
                    break;
                case 'Sunday':
                    dayIndex = 6;
                    break;

                default:
                    dayIndex = -1; // Invalid day
            }

            // Assuming you have an API endpoint for updating the repeat count
            const updateRepeatCountEndpoint = `/habit/updateRepeatCount/${habitId}`;

            // Fetch API to update the repeat count
            fetch(updateRepeatCountEndpoint, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    day_index: dayIndex,
                    completed: !target.classList.contains('completed'),  // Toggle the completed status
                }),
            })
            .then(response => response.json())
            .then(updatedHabit => {
                // Update the UI if needed
                console.log('Repeat count updated:', updatedHabit);

                // Toggle the 'clicked' class to update the cell style
                target.classList.toggle('clicked');
            })
            .catch(error => {
                console.error('Error updating repeat count:', error);
            });
        }
    });
});

document.addEventListener('DOMContentLoaded', function () {
    document.getElementById('delete-form').addEventListener('submit', function (event) {
    // Prevent the default form submission
    event.preventDefault();

    // Get the selected habit ID from the dropdown
    var habitId = document.getElementById('habit-select').value;

    // Construct the URL with the selected habit ID
    var url = '/habitdelete/' + habitId;

    // Submit the form with the constructed URL
    this.action = url;
    this.submit();
});
});


