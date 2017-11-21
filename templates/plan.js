function watchStartPlan() {
  $(document).on('keyup', '#plan_duration_input', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();

      const duration = $('#plan_duration_input').val();
      $.post('/ui/plan', JSON.stringify({ duration }), function(data) {
        $('#plan_duration').hide();
        $('#current_planning').html(`
          <div id="current_planning_header" class="flex flex-space-between">
            <span>Current planning:</span>
            <button id="dismiss_planning" class="btn btn-link">Dismiss</button>
          </div>
          ${data}
        `);
      });
    }
  });
}

function watchDismissPlan() {
  $(document).on('click', '#dismiss_planning', function(event) {
    $('#plan_duration').show();
    $('#current_planning').html('');
  });
}

$(document).ready(function() {
  watchStartPlan();
  watchDismissPlan();
});
