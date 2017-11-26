// Contains all the functions used to mark a task as done

// When clicking on a task, open the input to fill the description of
// how the task was done
function watchClickOnTasks(identifier) {
  $(identifier).on('click', '.TaskPending', function(event) {
    event.preventDefault();

    if ($(event.target).closest('.TaskDelete, .TaskDone, .TaskEdit', '#edit_input').length !== 0) {
      return;
    }

    if ($(this).find('#done_input').length) {
      return;
    }

    // Check if it is somewhere else
    if ($('#done_input').length) {
      $('#done_input').remove();
    }

    // Add it where it belongs
    $(this).find('.TaskDoneInputPlaceholder').html(`
      <textarea
        id="done_input"
        type="text"
        class="form-control"
        placeholder="How did you do that? (enter to mark as done)"
        data-taskid="${$(this).data('taskid')}"
      ></textarea>
    `);
    autosize($('#done_input'));
  });
}

// Watch all clicks outside a task row to hide a potential done input
function watchClickOutsideTask() {
  $(window).on('click', function(event) {
    if ($(event.target).closest('.TaskPending').length === 0) {
      $('#done_input').remove();

      // Check if it is somewhere else
      if ($('#edit_input').length) {
        $('#edit_input').remove();
        $('#edit_input_help').remove();
        $('.TaskContent').each(function(i, elt) {
          $(elt).show();
        });
      }
    }
  });
}

// Watch all the descendants of "identifier" matching ".TaskDone", marking the
// task done when such element is clicked
function watchDoneButtons(identifier) {
  $(identifier).on('click', '.TaskDone', function(event) {
    event.stopPropagation();

    $.post(
      `/ui/tasks/${$(this).data('taskid')}/done`,
      JSON.stringify({ description: $('#done_input').val() }),
      function(data) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);
        makeSortable();
        updateDoneTasks();
        refreshPlanning();
      },
    );
  });
}

// Mark a task as done when validating the done description
function watchDoneWithDescription(identifier) {
  $(identifier).on('keyup', '#done_input', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();
      $.post(
        `/ui/tasks/${$(this).data('taskid')}/done`,
        JSON.stringify({ description: $('#done_input').val() }),
        function(data) {
          $('#tasks_list ul').sortable('disable');
          $('#tasks_list').html(data);

          makeSortable();
          updateDoneTasks();
          refreshPlanning();
        },
      );
    } else if (event.keyCode === 27) {
      $('#done_input').remove();
    }
  });
}

function watchEdit() {
  $(document).on('click', '.TaskEdit', function(event) {
    event.preventDefault();

    const task = $(this).closest('.Task');

    if (task.find('#edit_input').length) {
      return;
    }

    // Check if it is somewhere else
    if ($('#edit_input').length) {
      $('#edit_input').remove();
      $('#edit_input_help').remove();
      $('.TaskContent').each(function(i, elt) {
        $(elt).show();
      });
    }

    task.find('.TaskContent').hide();
    task.find('.TaskEditPlaceholder').html(`
      <textarea id="edit_input" class="w-100 form-control no-border" data-taskid=${$(this).data(
        'taskid',
      )}>${task.find('.TaskEditPlaceholder').data('raw')}</textarea>
        <small id="edit_input_help" class="grey">
          Press enter to create <i class="fa fa-level-down fa-rotate-90"></i>
        </small>
    `);
    autosize($('#edit_input'));
  });
}

// Mark a task as done when validating the done description
function watchEditFinished() {
  $(document).on('keydown', '#edit_input', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();

      $.post(`/ui/tasks/${$(this).data('taskid')}`, JSON.stringify({ content: $('#edit_input').val() }), function(
        data,
      ) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);

        $('#edit_input').remove();
        $('#edit_input_help').remove();

        makeSortable();
        updateDoneTasks();
        refreshPlanning();
      });
    } else if (event.keyCode === 27) {
      $('#edit_input').remove();
      $('#edit_input_help').remove();
      $('.TaskContent').each(function(i, elt) {
        $(elt).show();
      });
    }
  });
}

$(document).ready(function() {
  watchEdit();
  watchDoneButtons(document);
  watchClickOnTasks(document);
  watchDoneWithDescription(document);
  watchClickOutsideTask();
  watchEditFinished();

  $(function() {
    $('[data-toggle="tooltip"]').tooltip();
  });
});
