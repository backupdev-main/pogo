program simple;

var a, b : int;

func simpleFunc(v : int, g : int) {
    print(v, g)
};

begin
    a = 5;
    b = 1;
    print("Starting nested loop test", a, b)

    while (a > 0) {
        print("Outer loop: a =", a)
        while (b < 3) {
            print("Inner loop: b =", b)
            b = b + 1;
        }
        a = a - 1;
        b = 1;  // Reset b for next iteration of outer loop
    }

    print("Loop finished. Final values: a =", a, "b =", b)
end