<?php
/* @var $this yii\web\View */
$this->title = 'Library index';
?>
<div class="site-index">

    <table cellpadding="2" cellspacing="2" border="0">
    <tr>
        <th>Created at</th>
        <th>Size (bytes)</th>
        <th>File path</th>
    </tr>

    <?php

        if (!empty($list['data'])) {
            $data = json_decode($list['data'], true);

            foreach ($data as $row) {
                echo "<tr><td>".$row['CreatedAt']."</td><td>".$row['SizeBytes']."</td>";
                echo "<td><a href='/site/info/?name=".basename($row['title'])."'>".basename($row['title'])."</a></td></tr>";
            }
        }

    ?>

    </table>

</div>
