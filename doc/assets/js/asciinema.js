import * as AsciinemaPlayer from 'asciinema-player';

AsciinemaPlayer.create('/asciinema/example.cast', document.getElementById("asciinema-player"), {cols: 200, rows: 30, fit: false, autoPlay: true});
