function watchStartPlan() {
  $(document).on('keyup', '#plan_duration_input', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();

      const input = $('#plan_duration_input').val();
      $.post('/ui/plan', JSON.stringify({ input }), function(data) {
        $('#current_planning').html(data);

        $(function() {
          $('[data-toggle="tooltip"]').tooltip();
        });
      }).fail(handleError);
    }
  });
}

function watchDismissPlan() {
  $(document).on('click', '#dismiss_planning', function(event) {
    $.ajax({
      url: '/ui/plan',
      type: 'DELETE',
      success: function(data) {
        $('#current_planning').html(data);
      },
    }).fail(handleError);
  });
}

function watchDoLater() {
  $(document).on('click', '.PlanningDoLater', function(event) {
    event.preventDefault();

    const taskId = $(this).data('taskid');
    $.post('/ui/plan/later', JSON.stringify({ taskId }), function(data) {
      $('#current_planning').html(data);

      $(function() {
        $('[data-toggle="tooltip"]').tooltip();
      });
    }).fail(handleError);
  });
}

function refreshPlanning() {
  $.get('/ui/plan', function(data) {
    $('#current_planning').html(data);

    $(function() {
      $('[data-toggle="tooltip"]').tooltip();
    });
  }).fail(handleError);
}

$(document).ready(function() {
  watchStartPlan();
  watchDismissPlan();
  watchDoLater();
});
