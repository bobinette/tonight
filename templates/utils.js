function realoadTooltips() {
  $(function() {
    $('[data-toggle="tooltip"]').tooltip();
  });
}

function notifyDanger(msg) {
  $.notify(
    {
      title: '<strong>Error</strong>',
      message: msg,
    },
    {
      type: 'danger',
      newest_on_top: true,
      animate: {
        enter: 'animated fadeIn',
        exit: 'animated fadeOut',
      },
    },
  );
}

function handleError(jqXHR, textStatus, errorThrown) {
  let msg = jqXHR.status;

  if (jqXHR.readyState === 0) {
    // Network error
    msg = 'Network issue';
  } else if (jqXHR.readyState === 4) {
    // HTTP error
    if (jqXHR.responseJSON.hasOwnProperty('error')) {
      msg = jqXHR.responseJSON.error;
    } else {
      msg = jqXHR.statusText;
    }
  }

  notifyDanger(msg);
}
