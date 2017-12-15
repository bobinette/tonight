function sortNumber(a, b) {
  return a - b;
}

function makeSortable() {
  $('#tasks_list ul').sortable({
    // Taken from http://www.guillaumevoisin.fr/jquery/tutoriel-drag-and-drop-jquery-exemple-avec-une-liste-des-taches
    axis: 'y', // Only on vertical axix
    containment: '#tasks_list', // Bound to this element
    handle: '.TaskDrag', // Only the .TaskDrag icon handles the drag
    distance: 10, // Drag starts after 10 pixel

    // When dropping
    stop: function(event) {
      // Pour chaque item de liste
      const availableRanks = [];
      const ranks = {};
      $('#tasks_list ul')
        .find('li')
        .each(function(idx, elt) {
          if (!$(elt).data('taskid')) {
            return;
          }

          availableRanks.push($(elt).data('taskrank'));
          ranks[$(elt).data('taskid')] = $(elt).index();
        });

      availableRanks.sort(sortNumber);
      const actualRanks = {};
      for (key in ranks) {
        actualRanks[key] = availableRanks[ranks[key]];
      }

      $.post(`/ui/ranks?q=${searchQ || ''}`, JSON.stringify({ ranks: actualRanks }), function(data) {
        $('#tasks_list ul').sortable('disable');
        $('#tasks_list').html(data);

        $(function() {
          $('[data-toggle="tooltip"]').tooltip();
        });

        makeSortable();
      }).fail(handleError);
    },
  });
}

$(document).ready(function() {
  makeSortable();
});
