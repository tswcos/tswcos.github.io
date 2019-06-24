function tableRemoveCounter(e) {
    if (e.checked) {
        document.getElementById("table_row_counter").disabled = true;
    } else {
        document.getElementById("table_row_counter").disabled = false;
    }
}