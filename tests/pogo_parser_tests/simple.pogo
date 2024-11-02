program simple;

var x, y: float;

func patito(x : int, y : int) {
    var z : int;
    // print(x + y + z)
};

begin
    x = (x + 5) * y + 7 / 5;
    if (x > 8) {
        y = x;
    } else {
        x = y;
        while (x > y) {
            x = x + 1;
        }
    }
    x = 55 + 95 / 9 * (5 * 4);
end