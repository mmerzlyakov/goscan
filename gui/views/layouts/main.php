<?php
$csrfParam = Yii::$app->request->csrfParam;
$csrfToken = Yii::$app->request->csrfToken;

use yii\helpers\Html;
use yii\widgets\Menu;

/* @var $this \yii\web\View */
/* @var $content string */

\yii\web\YiiAsset::register($this);
?>
<?php $this->beginPage() ?>
<!DOCTYPE html>
<html lang="<?= Yii::$app->language ?>">
<head>
    <meta charset="<?= Yii::$app->charset ?>">
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<?= Html::csrfMetaTags() ?>
    <title><?= Html::encode($this->title) ?></title>
    <link rel="stylesheet" href="<?= Yii::$app->request->getBaseUrl() ?>/css/site.css"/>
    <?php $this->head() ?>
</head>
<body>
<h1>Library app</h1>
<?php $this->beginBody() ?>
    <div class="header">
    <?= Menu::widget([
        'items' => [
            ['label' => 'Library', 'url' => ['/site/index']],
            Yii::$app->user->isGuest ? (
                ['label' => 'Login', 'url' => ['/site/login']]
            ) : (
                [
                    'url' => ['/site/logout'],
                    'label' => 'Logout (' . Yii::$app->user->identity->username . ')',
                    'template' => 

<<<HTML
<form method="post" action="{url}">
<input type="hidden" name="{$csrfParam}" value="{$csrfToken}" />
<input type="submit" value="{label}" />
</form>
HTML
,
                ]
            ),
        ]
    ]) ?>
    </div>

    <div class="content">
        <?= $content ?>
    </div>

    <footer class="footer">
        &copy; Library <?= date('Y') ?>, <?= Yii::powered() ?>
    </footer>
<?php $this->endBody() ?>
</body>
</html>
<?php $this->endPage() ?>
