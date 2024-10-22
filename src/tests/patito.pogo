program example123;

var var1, var2, var3 : int; //This is the only part where vars can be declared

func functionExample() {
    if (var1 > var2) {
        print(var1);
    };

    while (var1 > var2) {
        print(var1);
        var1 = var - 1;
    };
};

func functionExample2(param1 : int, param2 : float) {
    print(param1, param2);
};

func functionExample2(var3 : int, var4: float) {
    print(var3, var4);
};

begin
    functionExample();
    print("Example Worked");
end
