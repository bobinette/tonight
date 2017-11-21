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

$(document).ready(function() {
  watchNewTaskInput('#new_task_input');
});
