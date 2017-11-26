function watchStartPlan() {
  $(document).on('keyup', '#plan_duration_input', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();

      const duration = $('#plan_duration_input').val();
      $.post('/ui/plan', JSON.stringify({ duration }), function(data) {
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
});
