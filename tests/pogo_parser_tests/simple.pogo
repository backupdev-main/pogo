program ComprehensiveTest;

var x, y, z : int;
var a, b, c : float;

func testInt(x : int) {
    var result : int;
    print(x)
    print(x * 2 + 5 + 10)
    result = x + 8;
    print(result)
};

func testFloat(f : float, z : float) {
    var result : float;
    print(5.0/8)
    result = f * z + 5 / 8;
    print(result)
};

begin
    a = 5.4;
    b = 1.0;
    c = 2.5;

    x = 1;
    y = x+1;

    testInt(x)
    testInt(8 + 5)

    testFloat(b, c)

    if (x > y) {
        print("yupi")
    } else {
        print("not yupi")
        while (x < y) {
            x = x + 1;
            print(x)
        }
    }

    y = 5 / 8;

end