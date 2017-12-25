let showDoneTasks = false;

function updateDoneTasks() {
  if (showDoneTasks) {
    $.get('/ui/done', function(data) {
      $('#done_tasks_list').html(data);

      $(function() {
        $('[data-toggle="tooltip"]').tooltip();
      });
    }).fail(handleError);
  } else {
    $('#done_tasks_list').html('');
  }
}

function watchShowDoneTasks(identifer) {
  $(identifer).on('click', '#show_done_tasks_btn', function(event) {
    showDoneTasks = !showDoneTasks;
    updateDoneTasks();

    $('#show_done_tasks_btn').text(showDoneTasks ? 'Hide done tasks...' : 'Show done tasks...');
  });
}

$(document).ready(function() {
  watchShowDoneTasks(document);
});
