function makeSortable() {
  $('#tasks_list ul').sortable({
    // Taken from http://www.guillaumevoisin.fr/jquery/tutoriel-drag-and-drop-jquery-exemple-avec-une-liste-des-taches
    axis: 'y', // Only on vertical axix
    containment: '#tasks_list', // Bound to this element
    handle: '.TaskDrag', // Only the .TaskDrag icon handles the drag
    distance: 10, // Drag starts after 10 pixel

    // When dropping
    stop: function(event, ui) {
      // Pour chaque item de liste
      const ranks = {};
      $('#tasks_list ul')
        .find('li')
        .each(function() {
          ranks[$(this).data('taskid')] = $(this).index();
        });

      $.post('/ui/ranks', JSON.stringify({ ranks: ranks }), function(data) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);

        $(function() {
          $('[data-toggle="tooltip"]').tooltip();
        });

        makeSortable();
      });
    },
  });
}

$(document).ready(function() {
  makeSortable();
});
