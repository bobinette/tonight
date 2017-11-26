function watchSearchInput() {
  $('#search_input').on('keydown', function(event) {
    if (event.keyCode == 13) {
      event.preventDefault();

      $.get(`/ui/tasks?q=${$('#search_input').val()}`, function(data) {
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
