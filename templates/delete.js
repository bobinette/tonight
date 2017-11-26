function watchDeleteButtons(identifier) {
  $(identifier).on('click', '.TaskDelete', function(event) {
    event.preventDefault();

    $.ajax({
      url: `/ui/tasks/${$(this).data('taskid')}`,
      type: 'DELETE',
      success: function(data) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);
        makeSortable();

        $(function() {
          $('[data-toggle="tooltip"]').tooltip();
        });
      },
    }).fail(handleError);
  });
}

$(document).ready(function() {
  watchDeleteButtons(document);
});
