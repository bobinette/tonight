$.fn.visible = function() {
  return this.css('visibility', 'visible');
};

$.fn.invisible = function() {
  return this.css('visibility', 'hidden');
};

function watchNewTaskInput(identifer) {
  $(document).on('keyup', '#new_task_input', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();

      $.post(`/ui/tasks?q=${searchQ || ''}`, JSON.stringify({ content: $('#new_task_input').val() }))
        .done(function(data) {
          $('#new_task_input').val('');
          $('#tasks_list ul').sortable('disable');
          $('#tasks_list').html(data);

          if ($('#tasks_list').find('li').length > 0) {
            $('#new_task_input').addClass('HasTasks');
          } else {
            $('#new_task_input').removeClass('HasTasks');
          }

          makeSortable();

          $(function() {
            $('[data-toggle="tooltip"]').tooltip();
          });
        })
        .fail(handleError);
    }
  });
}

function watchAddTaskButton() {
  $(document).on('click', '#add_task_button', function(event) {
    event.preventDefault();
    $('#add_task_input').show();
    $('#add_task_input_textarea').focus();
  });
}

function watchClickOutsideAddTaskButton() {
  $(window).on('click', function(e) {
    if ($('#add_task_input').length) {
      if ($(e.target).closest('#add_task_input').length === 0 && $(e.target).closest('#add_task_button').length === 0) {
        // close/animate your div
        $('#add_task_input').hide();
      }
    }
  });
}

function watchAddTaskInput() {
  $('#add_task_input_textarea').on('keydown', function(event) {
    if (event.keyCode == 13) {
      event.preventDefault();

      $.post(`/ui/tasks?q=${searchQ || ''}`, JSON.stringify({ content: $('#add_task_input_textarea').val() }), function(
        data,
      ) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);

        if ($('#tasks_list').find('li').length > 0) {
          $('#new_task_input').addClass('HasTasks');
        } else {
          $('#new_task_input').removeClass('HasTasks');
        }

        makeSortable();

        $('#add_task_input_textarea').val('');
        $('#add_task_input_help').invisible();

        $.notify(
          {
            title: '<strong>Success</strong>',
            message: 'The task was successfully created.',
          },
          {
            type: 'success',
            newest_on_top: true,
            animate: {
              enter: 'animated fadeIn',
              exit: 'animated fadeOut',
            },
          },
        );

        $(function() {
          $('[data-toggle="tooltip"]').tooltip();
        });
      }).fail(handleError);
    } else if (event.keyCode == 27) {
      $('#add_task_input').hide();
    }
  });

  $('#add_task_input_textarea').on('keyup', function(event) {
    if ($('#add_task_input_textarea').val().length) {
      $('#add_task_input_help').visible();
    } else {
      $('#add_task_input_help').invisible();
    }
  });
}

$(document).ready(function() {
  autosize($('#add_task_input_textarea'));
  autosize($('#new_task_input'));

  watchNewTaskInput();

  watchAddTaskButton();
  watchClickOutsideAddTaskButton();
  watchAddTaskInput();
});
