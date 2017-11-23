$.fn.visible = function() {
  return this.css('visibility', 'visible');
};

$.fn.invisible = function() {
  return this.css('visibility', 'hidden');
};

function watchNewTaskInput(identifer) {
  $(identifer).on('keyup', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();

      $.post('/ui/tasks', JSON.stringify({ content: $('#new_task_input').val() }), function(data) {
        $('#new_task_input').val('');
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);
        makeSortable();
      });
    }
  });
}

function watchAddTaskButton() {
  $(document).on('click', '#add_task_button', function(event) {
    event.preventDefault();
    $('#add_task_input').show();
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

      $.post('/ui/tasks', JSON.stringify({ content: $('#add_task_input_textarea').val() }), function(data) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);
        makeSortable();

        $('#add_task_input').hide();
        $('#add_task_input_textarea').val('');
      });
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

  watchNewTaskInput('#new_task_input');

  watchAddTaskButton();
  watchClickOutsideAddTaskButton();
  watchAddTaskInput();
});