program simple;
var output : float;
var a, b, c, d, e, f, k, h, j, g, l: float;
var age_pablo, age_nan  : int;

func patito(x : int) {
    a = a + 3;
    print(5)
};

func patitow(x : float, y : float, z : int) {
    x = a + b / c * d * (-5 + 10 / 35);
    if (x > y) {
        while(x < y) {
            x = x + 1;
        }
    } else {
        patito(z)
        print("random print")
    }
};

begin
    print("Hello World", "Its me", a + b)
    a = a + c * 5 + (8 * 16 / 9 + 4);
    patito(age_pablo)

    print(age_pablo + age_nan)

    while (age_pablo * 5 > age_nan + 10 / 10 * 75) {
        print("Hola")
    }
end