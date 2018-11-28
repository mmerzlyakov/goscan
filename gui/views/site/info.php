<?php
/* @var $this yii\web\View */
$this->title = 'Details';

$data = json_decode($info['data'], true);

?>
<div class="site-index">

    <h1>Detailed info: <?= $data['title'] ?></h1>
    <table cellpadding="2" cellspacing="2" border="0">
        <tr>
            <th>Attribute</th>
            <th>Value</th>
        </tr>

        <tr>
            <td>File path</td>
            <td><input type="text" disabled id="ucount" value="<?= $data['title'] ?>"</td>
        </tr>
        <tr>
            <td>File name</td>
            <td><input type="text" disabled id="ucount" value="<?= basename($data['title']) ?>"</td>
        </tr>
        <tr>
            <td>File size (bytes)</td>
            <td><input type="text" disabled id="ucount" value="<?= $data['SizeBytes'] ?>"</td>
        </tr>
        <tr>
            <td>Words count</td>
            <td><input type="text" disabled id="ucount" value="<?= $data['wcount'] ?>"</td>
        </tr>
        <tr>
            <td>Unique words count</td>
            <td><input type="text" disabled id="ucount" value="<?= $data['ucount'] ?>"</td>
        </tr>
    </table>

    <h3>Check frequency for a word (case sensetive):  </h3>
    <input type="hidden" id="title" value="<?= basename($data['title']) ?>"></input>
    <input type="text" id="query" placeholder="type a word here..."></input>
    <input type="button" id="ask" value="Fetch"></input><br><br>
    Count of this word in text <input type="text" id="result" disabled></input><br><br>
    Frequency is <input type="text" id="frequncy" disabled></input> % percents

</div>

<script src="//ajax.googleapis.com/ajax/libs/jquery/3.1.0/jquery.min.js"></script>

<script>

    $("#ask").on('click',function () {
        var word = $("#query").val();
        var title = $("#title").val();
        $.ajax({
            type: 'GET',
            url: 'http://localhost:5051/api/stat/'+title+'/'+word,
            async: true,
            success: function (response) {
                var frequncy = (response.counter * 100)/response.wcount; //in percents
                console.log(frequncy);
                $('#result').val(response.counter);
                $('#frequncy').val(frequncy.toFixed(2));
            }
        });

    });

</script>
