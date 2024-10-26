program nestedTest;

var x : int;
var y, x : float;

begin
    if (x > y) {
        if (x == y) {
            print(x);
        } else {
            print(y);
        }
    }

    while (x > y) {
        x = x - 1;
        if (x != y) {
            print(x);
        }
    }
end