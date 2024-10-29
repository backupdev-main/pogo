program pablo;

var x, a : int;
var y, z : float;

func patito() {
    if (5 > 7) {
        x = -5;
        z = 7 + (x + y);
        x = 9;
    }
};

begin
    x = 5;
    a = 8 + ((5 + 40) + 8);
    if (x > y) {
        x = 5;
        z = 7 * 8.0;
        z = 7.0 / 8.0;
    }
    // Comments that should be ignored
    print("hola")

    while (x > y) {
        x = x + 1;
        if (x > 67) {
            print("wow")
        } else {
            print("no wow")
            patito()
            x = x * 2;
        }
    }

end