DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM exercise_library LIMIT 1) THEN
            INSERT INTO exercise_library (name, muscles, equipment, level) VALUES
                                                                               ('Martwy ciąg ze sztangą', 'Czworogłowe, Pośladki, Mięśnie dwugłowe uda, Pachwiny', 'Sztanga', 'Średniozaawansowany'),
                                                                               ('Martwy ciąg z kettlebell', 'Czworogłowe, Mięśnie dwugłowe uda, Pachwiny, Pośladki', 'Kettlebell', 'Początkujący'),
                                                                               ('Rumuński martwy ciąg z hantlami', 'Mięśnie dwugłowe uda, Pachwiny, Pośladki', 'Hantle', 'Średniozaawansowany'),
                                                                               ('Leżące uginanie nóg', 'Mięśnie dwugłowe uda, Łydki', 'Maszyna do uginania nóg na leżąco', 'Początkujący'),
                                                                               ('Uginanie ramion ze sztangą', 'Bicepsy', 'Sztanga', 'Początkujący'),
                                                                               ('Mostek biodrowy na podłodze', 'Mięśnie core, Pośladki, Pachwiny', 'Brak', 'Początkujący'),
                                                                               ('Kobra na brzuchu', 'Górna część pleców, Bark', 'Brak', 'Początkujący'),
                                                                               ('Skłony Good Mornings', 'Mięśnie dwugłowe uda, Pachwiny, Pośladki', 'Sztanga', 'Średniozaawansowany'),
                                                                               ('Deska boczna', 'Mięśnie brzucha', 'Brak', 'Początkujący'),
                                                                               ('Deska', 'Mięśnie brzucha', 'Brak', 'Początkujący'),
                                                                               ('Deska z chodzeniem', 'Mięśnie brzucha', 'Brak', 'Średniozaawansowany'),
                                                                               ('Deska na prostych ramionach', 'Mięśnie brzucha', 'Brak', 'Początkujący'),
                                                                               ('Przysiad więzienny', 'Uda, Pośladki, Łydki', 'Brak', 'Początkujący'),
                                                                               ('Wyskok z przysiadu', 'Całe ciało', 'Brak', 'Średniozaawansowany'),
                                                                               ('Bułgarski przysiad', 'Czworogłowe, Pośladki, Pachwiny', 'Ławka, Hantle', 'Średniozaawansowany'),
                                                                               ('Przysiad z hantlem (Goblet Squat)', 'Czworogłowe, Pachwiny, Pośladki', 'Kettlebell', 'Początkujący'),
                                                                               ('Przysiad z hantlami z przodu', 'Czworogłowe, Pachwiny, Pośladki', 'Hantle', 'Początkujący'),
                                                                               ('Burpees', 'Całe ciało', 'Brak', 'Zaawansowany'),
                                                                               ('Przysiad przedni z kettlebell', 'Uda, Czworogłowe, Pachwiny, Pośladki', 'Kettlebell', 'Początkujący'),
                                                                               ('Bird Dog', 'Mięśnie core, Prostownik grzbietu, Multifidi', 'Brak', 'Początkujący'),
                                                                               ('Podciąganie', 'Plecy, Mięsień najszerszy grzbietu, Romboidy, Trapezy, Bark, Tylne części barków, Ramiona, Brachialis', 'Drążek do podciągania', 'Średniozaawansowany'),
                                                                               ('Podciąganie z gumą', 'Plecy, Mięsień najszerszy grzbietu, Romboidy, Trapezy, Bark, Tylne części barków, Ramiona, Brachialis', 'Drążek, Guma lub taśma', 'Średniozaawansowany'),
                                                                               ('Rumuński martwy ciąg (sztanga)', 'Uda, Mięśnie dwugłowe uda, Pachwiny, Pośladki', 'Sztanga, Obciążenia', 'Średniozaawansowany'),
                                                                               ('Rosyjski twist', 'Mięśnie core, Skosy, Pośladki', 'Piłka stabilizacyjna', 'Średniozaawansowany'),
                                                                               ('Odwrotne spięcia brzucha z rotacją kolan', 'Mięśnie core, Prosty brzucha, Skosy, Plecy, Mięsień najszerszy grzbietu', 'Ławka', 'Średniozaawansowany'),
                                                                               ('Pompka', 'Klatka piersiowa, Bark, Przednie części barków, Tricepsy', 'Brak', 'Początkujący'),
                                                                               ('Pompka plyometryczna', 'Klatka piersiowa, Bark, Przednie części barków, Tricepsy', 'Brak', 'Średniozaawansowany'),
                                                                               ('Pompka pike', 'Barki, Przednie części barków, Barki, Trapezy, Tricepsy', 'Brak', 'Średniozaawansowany'),
                                                                               ('Pompka odwrócona', 'Barki, Przednie części barków, Barki, Trapezy, Tricepsy', 'Skrzynia lub stopień', 'Średniozaawansowany'),
                                                                               ('Pompka na ławce ujemnej', 'Klatka piersiowa, Bark, Przednie części barków, Tricepsy', 'Ławka', 'Średniozaawansowany'),
                                                                               ('Pompka na ławce dodatniej', 'Klatka piersiowa, Bark, Przednie części barków, Tricepsy', 'Ławka', 'Początkujący'),
                                                                               ('Pompka Archer', 'Klatka piersiowa, Bark, Przednie części barków, Tricepsy', 'Piłka medyczna', 'Zaawansowany'),
                                                                               ('Zmodyfikowana pompka', 'Klatka piersiowa, Bark, Przednie części barków, Tricepsy', 'Brak', 'Początkujący'),
                                                                               ('Wspięcia na palce na suwnicy', 'Łydki', 'Maszyna suwnicy', 'Początkujący'),
                                                                               ('Pajacyki', 'Całe ciało', 'Brak', 'Początkujący'),
                                                                               ('Suwnica', 'Uda, Czworogłowe, Pośladki, Pachwiny', 'Maszyna suwnicy', 'Początkujący'),
                                                                               ('Suwnica jednonóż', 'Uda, Czworogłowe, Pośladki, Pachwiny', 'Maszyna suwnicy', 'Początkujący'),
                                                                               ('Skoki na skrzynię', 'Pośladki, Uda, Mięśnie dwugłowe uda, Czworogłowe, Łydki', 'Skrzynia lub stopień', 'Średniozaawansowany'),
                                                                               ('Krzyż żelazny', 'Plecy, Prostownik grzbietu, Pośladki, Średni, Klatka piersiowa', 'Brak', 'Średniozaawansowany'),
                                                                               ('Uginanie ramion z kettlebell z przysiadem', 'Ramiona, Brachialis, Bicepsy', 'Kettlebell', 'Początkujący'),
                                                                               ('Face Pull', 'Barki, Plecy, Ramiona', 'Maszyna wyciągu, Lina', 'Początkujący'),
                                                                               ('Rozpiętki na wyciągu', 'Klatka piersiowa, Bark, Przednie części barków', 'Maszyna wyciągu', 'Początkujący'),
                                                                               ('Wiosłowanie gumą stojąc', 'Górna część pleców, Środkowa część pleców', 'Guma lub taśma', 'Początkujący'),
                                                                               ('Wyciskanie wąskim chwytem na ławce', 'Barki, Tricepsy', 'Ławka', 'Początkujący'),
                                                                               ('Dead Bug', 'Mięśnie core', 'Brak', 'Początkujący'),
                                                                               ('Wyskoki w wykroku', 'Całe ciało, Pośladki, Uda, Łydki', 'Brak', 'Średniozaawansowany'),
                                                                               ('Siedzące wiosłowanie blisko', 'Plecy, Mięsień najszerszy grzbietu, Romboidy, Trapezy, Bark, Tylne części barków, Ramiona, Bicepsy', 'Maszyna do wiosłowania', 'Początkujący'),
                                                                               ('Pozycja dziecka', 'Barki, Tylne części barków, Plecy, Mięsień najszerszy grzbietu, Prostownik grzbietu, Piszczel', 'Brak', 'Początkujący'),
                                                                               ('Wyskoki tuck jump', 'Całe ciało, Pośladki, Uda, Łydki', 'Brak', 'Średniozaawansowany'),
                                                                               ('Rozpiętki kablowe stojąc', 'Klatka piersiowa, Bark, Przednie części barków', 'Maszyna wyciągu', 'Początkujący'),
                                                                               ('Dips na ławce', 'Tricepsy, Bark, Przednie części barków', 'Ławka', 'Średniozaawansowany'),
                                                                               ('Wyciskanie maszynowe', 'Klatka piersiowa, Bark, Przednie części barków, Tricepsy', 'Maszyna do wyciskania', 'Początkujący'),
                                                                               ('Wyciskanie hantlem jednorącz', 'Klatka piersiowa, Bark, Przednie części barków', 'Ławka, Hantle', 'Średniozaawansowany'),
                                                                               ('Wyciskanie hantlami z gumą', 'Klatka piersiowa, Bark, Przednie części barków', 'Ławka, Hantle, Guma lub taśma', 'Średniozaawansowany'),
                                                                               ('Wyciskanie hantlem jednorącz na skosie', 'Klatka piersiowa, Bark, Przednie części barków', 'Ławka, Hantle', 'Średniozaawansowany'),
                                                                               ('Wyciskanie hantlami na skosie', 'Klatka piersiowa, Bark, Przednie części barków', 'Ławka, Hantle', 'Początkujący'),
                                                                               ('Leżące uginanie nóg: dwie nogi faza koncentryczna, jedna ekscentryczna', 'Uda, Mięśnie dwugłowe uda, Łydki', 'Maszyna do uginania nóg na leżąco', 'Średniozaawansowany'),
                                                                               ('Leżące uginanie nóg: jedna noga', 'Uda, Mięśnie dwugłowe uda, Łydki', 'Maszyna do uginania nóg na leżąco', 'Początkujący'),
                                                                               ('Siedzące uginanie nóg', 'Mięśnie dwugłowe uda, Łydki', 'Maszyna do uginania nóg na siedząco', 'Początkujący'),
                                                                               ('Siedzące uginanie nóg jednonóż', 'Uda, Mięśnie dwugłowe uda, Łydki', 'Maszyna do uginania nóg na siedząco', 'Początkujący'),
                                                                               ('Przysiad jednonóż', 'Czworogłowe, Pachwiny, Pośladki', 'Brak', 'Zaawansowany'),
                                                                               ('Przysiad jednonóż z dotknięciem', 'Czworogłowe, Pośladki, Pachwiny', 'Brak', 'Zaawansowany'),
                                                                               ('Przysiad jednonóż do wiosłowania', 'Czworogłowe, Pachwiny, Pośladki, Mięśnie dwugłowe uda, Bicepsy, Środkowa część pleców, Tylne części barków', 'Maszyna wyciągu', 'Zaawansowany'),
                                                                               ('Wyciskanie sztangi na ławce', 'Klatka piersiowa', 'Ławka, Sztanga, Obciążenia, Zabezpieczenia', 'Początkujący'),
                                                                               ('Wyciskanie sztangi na ławce z gumą', 'Klatka piersiowa', 'Ławka, Sztanga, Obciążenia, Zabezpieczenia, Guma lub taśma', 'Zaawansowany'),
                                                                               ('Wyciskanie sztangi na ławce z łańcuchami', 'Klatka piersiowa', 'Ławka, Sztanga, Obciążenia, Zabezpieczenia, Łańcuchy', 'Zaawansowany'),
                                                                               ('Wyciskanie sztangi na ławce skośnej', 'Klatka piersiowa', 'Ławka, Sztanga, Obciążenia, Zabezpieczenia', 'Średniozaawansowany'),
                                                                               ('Wałkowanie pianą: Przywodziciele', 'Pachwiny', 'Wałek do masażu', 'Początkujący'),
                                                                               ('Wałkowanie pianą: Łydki', 'Pachwiny, Łydki', 'Wałek do masażu', 'Początkujący'),
                                                                               ('Wałkowanie pianą: Mięsień najszerszy grzbietu', 'Środkowa część pleców', 'Wałek do masażu', 'Początkujący'),
                                                                               ('Statyczne: Rozciąganie motyla', 'Pachwiny', 'Brak', 'Początkujący'),
                                                                               ('Statyczne: Rozciąganie mięśnia najszerszego piłką', 'Mięsień najszerszy grzbietu', 'Piłka stabilizacyjna', 'Początkujący'),
                                                                               ('Statyczne: Rozciąganie łydek na siedząco', 'Łydki', 'Pasek do rozciągania', 'Początkujący'),
                                                                               ('Statyczne: Rozciąganie przywodzicieli na stojąco', 'Pachwiny', 'Brak', 'Początkujący'),
                                                                               ('Stanie na jednej nodze z wyciągnięciem do przodu', 'Całe ciało', 'Brak', 'Początkujący');
        END IF;
    END
$$;