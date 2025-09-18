import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.IOException;


public class main {

    public static void main(String[] args) throws IOException {
        BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));

        String line;

        while ((line = reader.readLine()) != null) {
            System.out.println(line);
        }
    }
}