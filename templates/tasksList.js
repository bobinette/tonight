// Contains all the functions used to mark a task as done

// When clicking on a task, open the input to fill the description of
// how the task was done
function watchClickOnTasks(identifier) {
  $(identifier).on('click', '.TaskPending', function(event) {
    event.stopPropagation();

    if ($(this).find('#done_input').length) {
      return;
    }

    // Check if it is somewhere else
    if ($('#done_input').length) {
      $('#done_input').remove();
    }

    // Add it where it belongs
    $(this).find('.TaskDoneInputPlaceholder').html(`
      <input
        id="done_input"
        type="text"
        class="form-control"
        placeholder="How did you do that? (enter to mark as done)"
        data-taskid="${$(this).data('taskid')}"
      >
    `);
  });
}

// Watch all clicks outside a task row to hide a potential done input
function watchClickOutsideTask() {
  $(window).on('click', function(event) {
    $('#done_input').remove();
  });
}

// Watch all the descendants of "identifier" matching ".TaskDone", marking the
// task done when such element is clicked
function watchDoneButtons(identifier) {
  $(identifier).on('click', '.TaskDone', function(event) {
    event.stopPropagation();

    $.post(`/ui/tasks/${$(this).data('taskid')}`, JSON.stringify({ description: $('#done_input').val() }), function(
      data,
    ) {
      $('#tasks_list ul').sortable('disable');
      $('#tasks_list').html(data);
      makeSortable();
      updateDoneTasks();
      refreshPlanning();
    });
  });
}

// Mark a task as done when validating the done description
function watchDoneWithDescription(identifier) {
  $(identifier).on('keyup', '#done_input', function(event) {
    if (event.keyCode === 13) {
      event.preventDefault();
      $.post(`/ui/tasks/${$(this).data('taskid')}`, JSON.stringify({ description: $('#done_input').val() }), function(
        data,
      ) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);
        makeSortable();
        updateDoneTasks();
        refreshPlanning();
      });
    } else if (event.keyCode === 27) {
      $('#done_input').remove();
    }
  });
}

$(document).ready(function() {
  watchDoneButtons(document);
  watchClickOnTasks(document);
  watchDoneWithDescription(document);
  watchClickOutsideTask();
});
