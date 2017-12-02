let searchQ = '';

function watchSearchInput() {
  $('#search_input').on('keydown', function(event) {
    if (event.keyCode == 13) {
      event.preventDefault();

      searchQ = $('#search_input').val();
      $.get(`/ui/tasks?q=${searchQ || ''}`, function(data) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);
        makeSortable();

        $(function() {
          $('[data-toggle="tooltip"]').tooltip();
        });
      }).fail(handleError);
    }
  });
}

$(document).ready(function() {
  watchSearchInput();
});
