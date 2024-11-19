program ComprehensiveTest;

var x : int;

func factorial(x : int) {
    var result : int;
    result = 1;
    while(x > 0) {
        result = result * x;
        x = x - 1;
    }

    print(result)
};

func factorial2() {
    factorial(5)
};

begin
    factorial(5)
    factorial2()
    print("Hellow World")
    x = 5;
    while (x > 0) {
        print("This is x: ", x)
        x = x - 1;
    }
end