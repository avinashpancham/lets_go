document.getElementById('button').onclick = function() {
    location.href = '/shortener/' + document.getElementById('searchbox').value
}